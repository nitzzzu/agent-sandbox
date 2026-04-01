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
	"time"

	"github.com/agent-sandbox/agent-sandbox/pkg/activator"
	"github.com/agent-sandbox/agent-sandbox/pkg/config"
	"github.com/agent-sandbox/agent-sandbox/pkg/sandbox"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// ScalingDownOfIdleTimeout checks sandboxes for idle timeout and deletes them if necessary
func (s *Scaler) ScalingDownOfIdleTimeout() {
	sbs, err := s.controller.ListAll()
	if err != nil {
		klog.Error("Failed to list sandboxes for idle timeout scaling down: ", err)
		return
	}

	for _, sb := range sbs {
		// Skip if IdleTimeout is not configured (0 or negative)
		if sb.IdleTimeout <= 0 {
			klog.V(2).Infof("Sandbox %v IdleTimeout is %d, skipping idle timeout check", sb.Name, sb.IdleTimeout)
			continue
		}

		// Get the last request time from events
		lastRequestTime := s.activator.GetLastRequestTime(sb.Name)

		// If no LastRequestTime event found, use creation time as fallback
		if lastRequestTime == 0 {
			createTime := sb.CreatedAt
			lastRequestTime = createTime.Unix()
			klog.V(2).Infof("Sandbox %v has no LastRequestTime event, using creation time %v", sb.Name, createTime)
		}

		// Calculate idle time
		now := time.Now().Unix()
		idleTime := now - lastRequestTime
		idleTimeout := int64(sb.IdleTimeout)

		klog.V(2).Infof("Sandbox %v idle check: lastRequestTime=%d, now=%d, idleTime=%d, idleTimeout=%d",
			sb.Name, lastRequestTime, now, idleTime, idleTimeout)

		// Check if sandbox has been idle for longer than IdleTimeout
		if idleTime > idleTimeout {
			klog.Infof("Sandbox %s has been idle for %d seconds (threshold: %d seconds), triggering idle policy: %s",
				sb.Name, idleTime, idleTimeout, sb.IdlePolicy)

			// Execute idle policy
			if err := s.executeIdlePolicy(sb); err != nil {
				klog.Errorf("Failed to execute idle policy for sandbox %v, error %v", sb.Name, err)
				continue
			}

			// Record event
			r := activator.GetRecorder(s.rootCtx)
			obj := &v1.ReplicaSet{
				TypeMeta: v1meta.TypeMeta{
					Kind:       "ReplicaSet",
					APIVersion: "apps/v1",
				},
				ObjectMeta: v1meta.ObjectMeta{
					Name:      sb.Name,
					Namespace: config.Cfg.SandboxNamespace,
				},
			}
			r.Event(obj, corev1.EventTypeNormal, "ScaleDownIdleTimeout",
				"Sandbox scaled down due to idle timeout")

			klog.Infof("Scaled down sandbox %s due to idle timeout. IdleTime: %ds, IdleTimeout: %ds, IdlePolicy: %s",
				sb.Name, idleTime, idleTimeout, sb.IdlePolicy)
		}
	}
}

// executeIdlePolicy executes the idle policy for a sandbox
func (s *Scaler) executeIdlePolicy(sb *sandbox.Sandbox) error {
	switch sb.IdlePolicy {
	case "delete":
		// Delete the sandbox
		return s.controller.Delete(sb.Name)
	case "scaledown":
		// Scale down to 0 replicas (to be implemented if needed)
		klog.Infof("ScaleDown policy for sandbox %s - to be implemented", sb.Name)
		return nil
	default:
		// Default to delete
		return s.controller.Delete(sb.Name)
	}
}
