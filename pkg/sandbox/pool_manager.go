/*
 * Copyright 2025 The https://github.com/agent-sandbox/agent-sandbox Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sandbox

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/agent-sandbox/agent-sandbox/pkg/config"
	v1 "k8s.io/api/apps/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
)

// PoolManager manages sandbox pool replicasets
type PoolManager struct {
	client  kubernetes.Interface
	rootCtx context.Context
}

// NewPoolManager creates a new PoolManager instance
func NewPoolManager(ctx context.Context) *PoolManager {
	c := kubeclient.Get(ctx)
	return &PoolManager{
		client:  c,
		rootCtx: ctx,
	}
}

// AcquirePoolReplicaSet tries to acquire a pool replicaset for use
// It marks the pool replicaset as in-use and updates its metadata with actual sandbox data
// Returns error if concurrent update conflict occurs (optimistic locking)
func (pm *PoolManager) AcquirePoolReplicaSet(sb *Sandbox) (*v1.ReplicaSet, error) {
	// Try to find an available pool replicaset
	available, err := pm.listAvailablePoolReplicaSets(sb.TemplateObj)
	if err != nil {
		return nil, fmt.Errorf("failed to list available pool replicasets  %v", err)
	}

	// If none available, create and return
	if len(available) == 0 {
		var createdRS *v1.ReplicaSet
		//  create, if pool size is 0
		createdRS, err = pm.createReplicaSet(sb)
		if err != nil || createdRS == nil {
			return nil, fmt.Errorf("failed to create available pool replicasets %v", err)
		}
		return createdRS, nil
	}

	// available pool replicasets exist, random select one
	// and adapt pool instance to sandbox instance
	var poolRS *v1.ReplicaSet
	//random select one available pool replicaset
	idx := time.Now().UnixNano() % int64(len(available))
	poolRS = available[idx]
	klog.Infof("Acquired available pool replicasets count %v, select %s", len(available), poolRS.Name)

	acquiredRS, err := pm.adaptReplicasetToSandbox(poolRS, sb)
	if err != nil {
		klog.Warningf("failed to acquire pool replicaset %s: %v", poolRS.Name, err)
		return nil, fmt.Errorf("failed to acquire any pool replicaset: %v", err)
	}

	return acquiredRS, nil
}

// ListAvailablePoolReplicaSets lists available pool replicasets for a given template
func (pm *PoolManager) listAvailablePoolReplicaSets(template *config.Template) ([]*v1.ReplicaSet, error) {
	// Build selector: sbx-pool=true, sbx-pool-in-use=false, sbx-template=<template>
	selector := labels.Set{
		PoolLabel: "true",
		TPLLabel:  template.Name,
	}.AsSelector()

	rsList, err := pm.client.AppsV1().ReplicaSets(config.Cfg.SandboxNamespace).List(context.TODO(), v1meta.ListOptions{
		LabelSelector: selector.String(),
	})

	//TODO retry with times
	if err != nil {
		klog.Errorf("failed to list pool replicasets for template %s: %v", template, err)
		return nil, err
	}

	result := []*v1.ReplicaSet{}
	for i := range rsList.Items {
		rs := &rsList.Items[i]
		rsImg := rs.Spec.Template.Spec.Containers[0].Image

		// check image is same as template, in case some pool replicaset is created with old template image when template is updated
		// if not same, delete it and skip
		if rsImg != template.Image {
			klog.Warningf("deleting pool replicaset %s with outdated image %s for template %s", rs.Name, rsImg, template.Name)
			err := pm.client.AppsV1().ReplicaSets(config.Cfg.SandboxNamespace).Delete(context.TODO(), rs.Name, v1meta.DeleteOptions{})
			if err != nil {
				klog.Errorf("failed to delete pool replicaset %s with outdated image %s for template %s: %v", rs.Name, rsImg, template.Name, err)
			}
			continue
		}
		result = append(result, rs)
	}

	klog.V(2).Infof("found %d available pool replicasets for template %s", len(result), template)
	return result, nil
}

func (pm *PoolManager) adaptReplicasetToSandbox(rs *v1.ReplicaSet, sb *Sandbox) (*v1.ReplicaSet, error) {
	// Create a deep copy for update
	rsCopy := rs.DeepCopy()

	sbLabels := sb.GetLabels()
	rsLabels := rsCopy.GetLabels()

	// Retain the pool replicaset ID as sandbox ID
	sb.ID = rsLabels[IDLabel]
	sb.Name = rsCopy.Name
	sbLabels[IDLabel] = rsLabels[IDLabel]
	sb.SetLabels(sbLabels)

	// Reset as in-use and update labels
	rsLabels[UserLabel] = sb.User
	rsLabels[PoolLabel] = "false"
	rsLabels[TimeLabel] = strconv.FormatInt(time.Now().Unix(), 10)
	rsCopy.SetLabels(rsLabels)

	// Reset annotations with actual sandbox data
	raw, _ := json.Marshal(sb)
	anns := rsCopy.GetAnnotations()
	if anns == nil {
		anns = make(map[string]string)
	}
	anns["sandbox-data"] = string(raw)
	rsCopy.SetAnnotations(anns)

	// Update the ReplicaSet in Kubernetes
	_, err := pm.client.AppsV1().ReplicaSets(config.Cfg.SandboxNamespace).Update(context.TODO(), rsCopy, v1meta.UpdateOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to update pool replicaset %s: %v", rs.Name, err)
	}

	klog.V(2).Infof("acquired pool replicaset %s for sandbox %s", rs.Name, sb.Name)
	return rsCopy, nil
}

// createReplicaSet creates a specified number of pool replicasets for a template
// force indicates whether to create at least one replicaset when not existing ones before list return empty
func (pm *PoolManager) createReplicaSet(sb *Sandbox) (*v1.ReplicaSet, error) {
	klog.Infof("creating pool replicasets for sandbox %v", sb)

	var createdRS *v1.ReplicaSet

	var err error
	if createdRS, err = pm.createSingleReplicaSet(sb); err != nil {
		klog.Errorf("failed to create pool replicaset for sandbox %+v: %v", sb, err)
		return nil, err
	}

	return createdRS, nil
}

// createSingleReplicaSet creates a single pool replicaset
func (pm *PoolManager) createSingleReplicaSet(sb *Sandbox) (*v1.ReplicaSet, error) {
	// Create a pool sandbox with same configuration as user request
	rs, err := buildReplicaSet(sb)
	if err != nil {
		return nil, fmt.Errorf("failed to build replicaset from sandbox %v error %v", sb, err)
	}

	lbs := rs.GetLabels()
	lbs[PoolLabel] = "false"
	// Set pool labels
	if sb.IsPool {
		lbs[PoolLabel] = "true"
	}
	rs.SetLabels(lbs)

	// Create ReplicaSet in Kubernetes
	createdRS, err := pm.client.AppsV1().ReplicaSets(config.Cfg.SandboxNamespace).Create(context.TODO(), rs, v1meta.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create pool replicaset in kubernetes replicaset: %+v error: %+v", rs, err)
	}

	klog.V(2).Infof("created pool replicaset for sandbox %v", sb)
	return createdRS, nil
}

func (pm *PoolManager) StartPoolSyncing() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			pm.replenishPoolAsync()
		case <-pm.rootCtx.Done():
			klog.Info("Scaler stopping")
			return
		}
	}
}

// ReplenishPoolAsync asynchronously replenishes the pool to the target size
// TODO: delete pool, current delete pool by manually: kubectl -n sandbox delete rs -l sbx-pool=true sbx-template=<template-name>
func (pm *PoolManager) replenishPoolAsync() {
	tpls := config.Templates
	for _, tpl := range tpls {
		// skip if type is dynamic, since dynamic template image is determined by regexp pattern,
		// if what to use pool, define a static template with static image
		if tpl.Type == "dynamic" {
			continue
		}

		count := tpl.Pool.Size
		if count <= 0 {
			continue
		}

		// check exist pool rs
		sbs, err := pm.listAvailablePoolReplicaSets(tpl)
		if err != nil {
			klog.Errorf("failed to list available pool replicasets when CreatePoolReplicaSets %s, error: %v", tpl.Name, err)
			continue
		}

		// reduce exist pool rs
		count = count - len(sbs)

		if count <= 0 {
			continue
		}

		klog.V(2).Infof("creating pool replicasets for sandbox %v, count %v", tpl.Name, count)
		for i := 0; i < count; i++ {
			sb := GetDefaultSandbox()
			sb.Template = tpl.Name
			sb.IsPool = true
			sb.User = "sys-2492a85b10ed4cb083b2c76b181eac96"

			// init name and valid fields
			if err := sb.Make(); err != nil {
				klog.Errorf("failed to init pool sandbox for template %s: %v", tpl.Name, err)
				continue
			}

			_, err := pm.createReplicaSet(sb)
			if err != nil {
				klog.Errorf("failed to create pool replicaset %d/%d for sandbox %v: %v", i+1, count, sb, err)
			} else {
				klog.V(2).Infof("successfully replenished pool replicasets for template %s", tpl.Name)
			}
		}

	}
}
