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

package scaler

import (
	"context"
	"time"

	"github.com/agent-sandbox/agent-sandbox/pkg/activator"
	"github.com/agent-sandbox/agent-sandbox/pkg/sandbox"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type Scaler struct {
	rootCtx    context.Context
	activator  *activator.Activator
	controller *sandbox.Controller
}

func NewScaler(ctx context.Context, a *activator.Activator, c *sandbox.Controller) *Scaler {
	scaler := &Scaler{
		rootCtx:    ctx,
		activator:  a,
		controller: c,
	}
	return scaler
}

func (s *Scaler) RunScaling() {
	// Periodically check for sandboxes to scale down
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.ScalingDownOfTimeout()
		case <-s.rootCtx.Done():
			klog.Info("Scaler stopping")
			return
		}
	}
}

func (s *Scaler) ScalingDownOfTimeout() {
	sbs, err := s.controller.ListAll()
	if err != nil {
		klog.Error("Failed to list sandboxes for scaling down: ", err)
		return
	}

	for _, sb := range sbs {
		createT := sb.GetCreationTimestamp()
		timeout := sb.Timeout
		if timeout == 0 {
			klog.V(2).Infof("Sandbox %v timeout is 0, skipping scaling down", sb.Name)
			continue
		}
		tt := createT.Add(time.Duration(timeout) * time.Second)
		if tt.Before(time.Now()) {
			if err := s.controller.Delete(sb.Name); err != nil {
				klog.Errorf("Failed to scale down sandbox %v, error %v", sb, err)
				continue
			}
			r := activator.GetRecorder(s.rootCtx)
			obj := &v1.ReplicaSet{
				TypeMeta: v1meta.TypeMeta{
					Kind:       "ReplicaSet",
					APIVersion: "apps/v1",
				},
				ObjectMeta: v1meta.ObjectMeta{
					Name:      sb.Name,
					Namespace: sb.GetNamespace(),
				},
			}
			r.Event(obj, corev1.EventTypeNormal, "ScaleDownTimeout", "Sandbox scaled down due to timeout")
			klog.Infof("Scaled down sandbox %s CreationTimestamp %s Timeout %v IdleTimeout %v", sb.Name, sb.GetCreationTimestamp(), sb.Timeout, sb.IdleTimeout)
		}

	}

}
