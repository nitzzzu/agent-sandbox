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

package utils

import (
	"github.com/agent-sandbox/agent-sandbox/pkg/sandbox"
	resource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
)

type ResourceQuantity struct {
	CPUMilli   int64
	MemoryMB   int64
	DiskSizeMB int64
}

func CalculateResourceToQuantity(sb *sandbox.Sandbox) ResourceQuantity {
	cpu, err := resource.ParseQuantity(sb.CPU)
	if err != nil {
		klog.Errorf("failed to parse CPU quantity: %v", err)
		cpu = resource.MustParse("0")
	}

	mem, err := resource.ParseQuantity(sb.Memory)
	if err != nil {
		klog.Errorf("failed to parse Memory quantity: %v", err)
		mem = resource.MustParse("0")
	}

	rq := ResourceQuantity{
		CPUMilli:   cpu.MilliValue(),
		MemoryMB:   mem.Value() / (1024 * 1024),
		DiskSizeMB: 0,
	}

	return rq
}
