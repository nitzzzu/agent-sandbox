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

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	podclient "knative.dev/pkg/client/injection/kube/informers/core/v1/pod"

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
	tmpl, err := template.New(sb.Name).Parse(config.SandboxDeployTemplate)
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

	// set startupProbe port to pool probe port if is pool rs
	if sb.IsPool {
		container := rsObj.Spec.Template.Spec.Containers[0]
		if container.StartupProbe != nil {
			container.StartupProbe.TCPSocket.Port = intstr.FromInt(sb.TemplateObj.Pool.ProbePort)
		}
	}

	return rsObj, nil
}

// WaitForReplicaSetReady waits for a ReplicaSet to become ready
func (s *Controller) WaitForReplicaSetReady(sb *Sandbox) error {
	return wait.PollUntilContextTimeout(context.TODO(), 500*time.Millisecond, 1*time.Minute, true, func(ctx context.Context) (bool, error) {
		rsCreated, err := s.kclient.AppsV1().ReplicaSets(config.Cfg.SandboxNamespace).Get(context.TODO(), sb.Name, v1meta.GetOptions{})
		if err != nil {
			return false, err
		}
		// Check if the ReplicaSet is ready
		replicas := *rsCreated.Spec.Replicas
		if replicas == rsCreated.Status.ReadyReplicas {
			klog.Infof("ReplicaSet %s in namespace %s is ready. Desired: %d, Ready: %d",
				sb.Name, config.Cfg.SandboxNamespace, replicas, rsCreated.Status.ReadyReplicas)
			return true, nil
		} else {
			klog.V(2).Infof("ReplicaSet %s in namespace %s is NOT ready. Desired: %d, Ready: %d\n",
				sb.Name, config.Cfg.SandboxNamespace, replicas, rsCreated.Status.ReadyReplicas)
			return false, nil
		}
	})
}

// StartupPoolReplicaSet waits for a ReplicaSet to become ready by executing a command in the pod and checking the port is listening
func (s *Controller) StartupPoolReplicaSet(sb *Sandbox, checkReady bool) error {
	selector := labels.Set{
		"sandbox": sb.Name,
	}.AsSelector()
	podList, err := podclient.Get(s.rootCtx).Lister().Pods(config.Cfg.SandboxNamespace).List(selector)

	if err != nil {
		return err
	}
	if len(podList) == 0 {
		return fmt.Errorf("pod list is empty, name: %s", sb.Name)
	}
	pod := podList[0]

	// 1, start up services, services should be listening port that config in template !
	cmd := []string{
		"sh",
		"-c",
		sb.TemplateObj.Pool.StartupCmd,
	}
	klog.Infof("execute startup command in pod %s: %s", pod.Name, sb.TemplateObj.Pool.StartupCmd)
	_, _, err = s.ExecCommandInPod(pod.Name, cmd)
	if err != nil {
		return err
	}

	// if checkReady is false, skip checking port listening,
	// since some template startup.sh have ready check in it.
	if !checkReady {
		return nil
	}

	// 2. use curl port to check if the port is listening when services is started
	klog.Infof("check if port %d is listening in pod %s", sb.Port, pod.Name)
	return wait.PollUntilContextTimeout(context.TODO(), 300*time.Millisecond, 1*time.Minute, true, func(ctx context.Context) (bool, error) {
		cmd := []string{
			"sh",
			"-c",
			fmt.Sprintf("curl -s -o /dev/null  http://127.0.0.1:%d", sb.Port),
		}
		_, _, err = s.ExecCommandInPod(pod.Name, cmd)
		if err != nil {
			klog.Infof("port %d is not listening in pod %s, retry...", sb.Port, pod.Name)
			return false, nil
		}
		return true, nil
	})
}
