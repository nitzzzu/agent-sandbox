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
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"context"
	"text/template"

	"github.com/agent-sandbox/agent-sandbox/pkg/config"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"

	"sigs.k8s.io/yaml"

	v1 "k8s.io/api/apps/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// buildReplicaSet Build a Kubernetes ReplicaSet from a Sandbox object
func buildReplicaSet(sb *Sandbox) (*v1.ReplicaSet, error) {
	raw, _ := json.Marshal(sb)
	tplData := &SandboxKube{
		Sandbox:   sb,
		RawData:   string(raw),
		Namespace: config.Cfg.SandboxNamespace,
	}
	tmpl, err := template.New(sb.Name).Parse(SandboxDeployTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse template fail: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tplData)
	if err != nil {
		return nil, fmt.Errorf("execute template fail: %v", err)
	}

	rsObj := &v1.ReplicaSet{}
	if err = yaml.Unmarshal(buf.Bytes(), rsObj); err != nil {
		return nil, fmt.Errorf("unmarshal template fail: %v", err)
	}

	rsAnns := rsObj.GetAnnotations()
	rsLbs := rsObj.GetLabels()
	for k, v := range tplData.Sandbox.GetAnnotations() {
		rsAnns[k] = v
	}
	for k, v := range tplData.Sandbox.GetLabels() {
		rsLbs[k] = v
	}
	rsObj.SetAnnotations(rsAnns)
	rsObj.SetLabels(rsLbs)

	// Set annotations and labels to pod template
	podAnns := rsObj.Spec.Template.GetAnnotations()
	podLbs := rsObj.Spec.Template.GetLabels()
	if podAnns == nil {
		podAnns = make(map[string]string)
	}
	if podLbs == nil {
		podLbs = make(map[string]string)
	}
	for k, v := range rsAnns {
		podAnns[k] = v
	}
	for k, v := range rsLbs {
		podLbs[k] = v
	}
	rsObj.Spec.Template.SetAnnotations(podAnns)
	rsObj.Spec.Template.SetLabels(podLbs)

	return rsObj, nil
}

// WaitForReplicaSetReady waits for a ReplicaSet to become ready
func (s *Controller) WaitForReplicaSetReady(name string) error {
	return wait.PollUntilContextTimeout(context.TODO(), 500*time.Millisecond, 1*time.Minute, true, func(ctx context.Context) (bool, error) {
		rsCreated, err := s.client.AppsV1().ReplicaSets(config.Cfg.SandboxNamespace).Get(context.TODO(), name, v1meta.GetOptions{})
		if err != nil {
			return false, err
		}
		// Check if the ReplicaSet is ready
		replicas := *rsCreated.Spec.Replicas
		if replicas == rsCreated.Status.ReadyReplicas {
			klog.Infof("ReplicaSet %s in namespace %s is ready. Desired: %d, Ready: %d",
				name, config.Cfg.SandboxNamespace, replicas, rsCreated.Status.ReadyReplicas)
			return true, nil
		} else {
			klog.V(2).Infof("ReplicaSet %s in namespace %s is NOT ready. Desired: %d, Ready: %d\n",
				name, config.Cfg.SandboxNamespace, replicas, rsCreated.Status.ReadyReplicas)
			return false, nil
		}
	})
}
