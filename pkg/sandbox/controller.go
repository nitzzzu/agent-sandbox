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
	"encoding/json"
	"fmt"
	"sort"

	"context"

	"github.com/agent-sandbox/agent-sandbox/pkg/config"
	"github.com/agent-sandbox/agent-sandbox/pkg/utils"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	rsclient "knative.dev/pkg/client/injection/kube/informers/apps/v1/replicaset"
	podclient "knative.dev/pkg/client/injection/kube/informers/core/v1/pod"

	v1core "k8s.io/api/core/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Controller struct {
	kclient kubernetes.Interface
	kcfg    *rest.Config
	rootCtx context.Context
	pl      *PoolManager
}

func NewController(ctx context.Context, cfg *rest.Config, pl *PoolManager) *Controller {
	c := kubeclient.Get(ctx)
	sh := &Controller{
		rootCtx: ctx,
		kclient: c,
		kcfg:    cfg,
		pl:      pl,
	}
	return sh
}

func (s *Controller) GetRSByID(id string) (*v1.ReplicaSet, error) {
	selector := labels.Set{IDLabel: id}.AsSelector()
	rss, err := rsclient.Get(s.rootCtx).Lister().List(selector)
	if err != nil {
		klog.Errorf("Failed to list rs, id %s error %v", id, err)
		return nil, err
	}
	if len(rss) == 0 {
		klog.Warningf("No rs found with id %s", id)
		return nil, fmt.Errorf("no rs found with id %s", id)
	}
	return rss[0], nil
}

func (s *Controller) GetByID(id string) (*Sandbox, error) {
	rs, err := s.GetRSByID(id)
	if err != nil {
		return nil, err
	}
	return s.GetSandbox(rs)
}

func (s *Controller) Get(name string) (*Sandbox, error) {
	selector, _ := labels.Parse(fmt.Sprintf("sandbox=%s", name))
	rss, err := rsclient.Get(s.rootCtx).Lister().ReplicaSets(config.Cfg.SandboxNamespace).List(selector)
	if err != nil {
		return nil, err
	}

	if len(rss) == 0 {
		return nil, fmt.Errorf("no Sandbox found with name %s", name)
	}

	return s.GetSandbox(rss[0])
}

func (s *Controller) GetSandbox(rs *v1.ReplicaSet) (*Sandbox, error) {
	raw := rs.Annotations["sandbox-data"]
	sb := &Sandbox{}
	json.Unmarshal([]byte(raw), sb)
	sb.Object = rs.DeepCopy()

	// Set the status of the sandbox
	replicas := *rs.Spec.Replicas
	if replicas == rs.Status.ReadyReplicas {
		sb.Status = Running
	} else {
		sb.Status = Creating
	}

	return sb, nil
}

func (s *Controller) ListAll() ([]*Sandbox, error) {
	selector := labels.Set{
		"owner":   "agent-sandbox",
		PoolLabel: "false",
	}.AsSelector()

	return s.DoList(selector)
}

func (s *Controller) List(user string) ([]*Sandbox, error) {
	selector, _ := labels.Parse(fmt.Sprintf("%s=%s", UserLabel, user))
	return s.DoList(selector)
}

func (s *Controller) DoList(selector labels.Selector) ([]*Sandbox, error) {
	rss, err := rsclient.Get(s.rootCtx).Lister().List(selector)
	if err != nil {
		klog.Errorf("failed to list sandboxes: %v", err)
		return nil, err
	}
	var sandboxes = []*Sandbox{}
	for _, rs := range rss {
		raw := rs.Annotations["sandbox-data"]
		sb := &Sandbox{}
		json.Unmarshal([]byte(raw), sb)
		sb.Object = rs.DeepCopy()
		// Set the status of the sandbox
		replicas := *rs.Spec.Replicas
		if replicas == rs.Status.ReadyReplicas {
			sb.Status = Running
		} else {
			sb.Status = Creating
		}

		sandboxes = append(sandboxes, sb)
	}
	//sort by CreatedAt desc
	sort.Slice(sandboxes, func(i, j int) bool {
		return sandboxes[i].CreatedAt.After(sandboxes[j].CreatedAt)
	})
	return sandboxes, nil
}

func IsAcquireError(err error) bool {
	return err != nil
}

func (s *Controller) Create(sb *Sandbox) (*Sandbox, error) {
	// retry to AcquirePoolReplicaSet if error is conflict,
	//because multiple sandboxes may try to acquire the same pool replicaset
	acquired := &v1.ReplicaSet{}
	fromPool := false
	err := retry.OnError(retry.DefaultRetry, IsAcquireError, func() error {
		var err error
		acquired, fromPool, err = s.pl.AcquirePoolReplicaSet(sb)
		return err
	})

	if err != nil {
		klog.Errorf("failed to acquire pool replica set, reqeust %v, error %v", sb, err)
		return nil, err
	}

	sb.Object = acquired

	// Wait for ReplicaSet to be ready
	if fromPool && sb.TemplateObj.Pool.StartupCmd != "" {
		if perr := s.StartupAndWaitForPoolReplicaSetReady(sb, false); perr != nil {
			klog.Errorf("timeout waiting for sandbox from pool to be ready: %v, instance not ready yet, please get it leater or check pod status", perr)
			return sb, nil
		}
	} else {
		if perr := s.WaitForReplicaSetReady(sb); perr != nil {
			klog.Errorf("timeout waiting for sandbox to be ready: %v, instance not ready yet, please get it leater or check pod status", perr)
			return sb, nil
		}
	}

	sb.Status = Running
	return sb, nil
}

func (s *Controller) GetInstances(name string) []*v1core.Pod {
	selector, _ := labels.Parse(fmt.Sprintf("sandbox=%s", name))
	pods, err := podclient.Get(context.TODO()).Lister().List(selector)
	if err != nil {
		return nil
	}
	return pods
}

func (s *Controller) DeleteByID(id string) error {
	rs, err := s.GetRSByID(id)
	if err != nil {
		return err
	}
	return s.Delete(rs.Name)
}

func (s *Controller) Delete(name string) error {
	selector, _ := labels.Parse(fmt.Sprintf("sandbox=%s", name))
	// delete rs by selector, since rs name may be different when acquire from pool
	rss, err := rsclient.Get(s.rootCtx).Lister().ReplicaSets(config.Cfg.SandboxNamespace).List(selector)
	if err != nil {
		return err
	}
	for _, rs := range rss {
		err = s.kclient.AppsV1().ReplicaSets(config.Cfg.SandboxNamespace).Delete(context.TODO(), rs.Name, v1meta.DeleteOptions{})
		if err != nil {
			klog.Errorf("failed to delete replicaset %s: %v", rs.Name, err)
			return err
		}
		return err
	}

	return nil
}

func (s *Controller) ExecCommandInPod(name string, cmd []string) (output string, outputErr string, err error) {
	return utils.ExecCommand(s.kclient, s.kcfg, config.Cfg.SandboxNamespace, name, "sandbox", cmd)
}
