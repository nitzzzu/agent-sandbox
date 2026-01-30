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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/agent-sandbox/agent-sandbox/pkg/config"
)

var SandboxDeployTemplate string

func init() {
	SandboxDeployTemplate = defaultSandboxTemplate
	if config.Cfg.SandboxTemplateFile != "" {
		var err error
		var val []byte
		if val, err = os.ReadFile(config.Cfg.SandboxTemplateFile); err != nil {
			panic(err)
		}
		SandboxDeployTemplate = string(val)
	}
}

// Defines values for SandboxState.
type SandboxState string

const (
	Paused   SandboxState = "paused"
	Running  SandboxState = "running"
	Creating SandboxState = "creating"
	Ready    SandboxState = "ready"
	Unready  SandboxState = "unready"
)

type SandboxBase struct {

	// Optionally give the sandbox a name.
	Name string `json:"name,omitempty" required:"false" jsonschema:"The unique name of Sandbox."`

	// The type to run as the container for the sandbox when Image is not set. e.g. aio/python/shell/
	Template string `json:"template,omitempty" required:"false" jsonschema:"The sandbox Template name."`
}

type SandboxObject struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

func (in *SandboxObject) DeepCopyObject() runtime.Object {
	return in
}

type Sandbox struct {
	SandboxBase

	// For K8s object metadata access
	metav1.Object `json:"-"`

	// Optionally give the sandbox a unique id.  compatible with E2B API
	ID string `json:"id,omitempty" required:"false" jsonschema:"The unique id of Sandbox."`

	// Set the CMD of the SandboxHandler, overriding any CMD of the container image.
	Args []string `json:"args,omitempty"`

	// Associate the sandbox with an app. Required unless creating from a container.
	App string `json:"app,omitempty" jsonschema:"App to for associate the sandbox with an app"`

	// The image to run as the container for the sandbox.
	Image string `json:"image,omitempty"`

	// Environment variables to set in the SandboxHandler.
	EnvVars map[string]string `json:"envVars,omitempty"`

	// Maximum lifetime of the sandbox in seconds. timeout is reached apply to delete action
	Timeout int `json:"timeout,omitempty" default:"1800"` // default 30m, 0 is no timeout

	// The amount of time in seconds that a sandbox can be idle before being terminated.
	IdleTimeout int `json:"idle_timeout,omitempty"` // default 10m

	// Policy to apply when the idle is reached. Options are 'delete' or 'scaledown'.
	IdlePolicy string `json:"idle_policy,omitempty" default:"delete"` // default delete

	// Working directory of the sandbox.
	Workdir string `json:"workdir,omitempty"`

	// CPU request
	CPU string `json:"cpu,omitempty"  default:"100m"`

	// Memory request
	Memory string `json:"memory,omitempty"  default:"128Mi"`

	// CPU limit
	CPULimit string `json:"cpu_limit,omitempty"  default:"1000m"`

	// Memory limit
	MemoryLimit string `json:"memory_limit,omitempty"  default:"1024Mi"`

	// Port for startup probe and main service
	Port int `json:"port,omitempty"  default:"8080"`

	// Status of the sandbox. Options are 'creating', 'running', 'idle', 'deleting', 'error'.
	Status SandboxState `json:"status,omitempty"`

	// CreatedAt Time when the sandbox was started
	CreatedAt time.Time `json:"created_at,omitempty"`

	Metadata map[string]string `json:"metadata,omitempty"`
}

const (
	IDLabel   = "sbx-id"
	UserLabel = "sbx-user"
	TPLLabel  = "sbx-template"
)

func GetDefaultSandbox(user string) *Sandbox {
	id := uuid.NewString()
	// remove '-'
	id = strings.Replace(id, "-", "", -1)

	sb := &Sandbox{
		ID:          id,
		CPU:         "100m",
		Memory:      "128Mi",
		CPULimit:    "2000m",
		MemoryLimit: "4024Mi",
		Timeout:     30 * 60, // 30 minutes
		IdleTimeout: 0,       // no idle timeout
		IdlePolicy:  "delete",
		Port:        8080,
	}
	sb.Object = &SandboxObject{}
	sb.Object.SetAnnotations(map[string]string{})
	sb.Object.SetLabels(map[string]string{})

	labels := sb.GetLabels()
	labels[IDLabel] = sb.ID
	labels[UserLabel] = user
	sb.SetLabels(labels)

	return sb
}

func (o *Sandbox) Make() error {
	// one day max
	maxTimeout := 60 * 60 * 24
	if o.Timeout >= maxTimeout && o.Timeout < 9999 {
		o.Timeout = maxTimeout
	}

	// default timeout, since int default is 0 when not set, so we set it to 30 minutes, E2B default is 300s
	if o.Timeout == 0 {
		o.Timeout = 30 * 60 // default 30 minutes
	}

	// no timeout, when set to 9999
	if o.Timeout == 9999 {
		o.Timeout = 0 // no timeout
	}

	// one hour max
	if o.IdleTimeout > 60*60 {
		o.IdleTimeout = 60 * 60
	}

	if o.Template == "" && o.Image == "" {
		o.Image = config.Cfg.SandboxDefaultImage
		o.Template = config.Cfg.SandboxDefaultTemplate
	}

	if o.Template != "" && o.Image == "" {
		tpl, err := config.GetTemplateByName(o.Template)
		if err != nil {
			return fmt.Errorf("failed to get Template by name %s: %v", o.Template, err)
		}
		o.Image = tpl.Image
		if tpl.Port != 0 {
			o.Port = tpl.Port
		}
	}

	if o.Template == "" && o.Image != "" {
		o.Template = "custom"
	}

	if o.Name == "" {
		o.Name = fmt.Sprintf("sandbox-%s-%d", o.Template, time.Now().Unix())
	}

	labels := o.GetLabels()
	labels[TPLLabel] = o.Template
	o.SetLabels(labels)

	return nil
}

type SandboxKube struct {
	Sandbox   *Sandbox
	RawData   string
	Namespace string
}

const defaultSandboxTemplate = `apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: {{.Sandbox.Name}}
  namespace: {{.Namespace}}
  annotations:
    sandbox-data:  |
        {{.RawData}}
  labels:
    sandbox: "{{.Sandbox.Name}}"
    owner: agent-sandbox
spec:
  replicas: 1
  selector:
    matchLabels:
      sandbox: "{{.Sandbox.Name}}"
  template:
    metadata:
      labels:
        sandbox: "{{.Sandbox.Name}}"
        owner: agent-sandbox
    spec:
      containers:
      - name: sandbox
        image: {{.Sandbox.Image}}
        imagePullPolicy: IfNotPresent
        env:
        - name: INSTANCE_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        resources:
          requests:
            cpu: {{.Sandbox.CPU}}
            memory: {{.Sandbox.Memory}}
          limits:
            cpu: {{.Sandbox.CPULimit}}
            memory: {{.Sandbox.MemoryLimit}}
        startupProbe:
            failureThreshold: 600
            tcpSocket:
              port: {{.Sandbox.Port}}
            periodSeconds: 1
            successThreshold: 1
            timeoutSeconds: 3

`
