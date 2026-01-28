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

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agent-sandbox/agent-sandbox/pkg/config"
	"github.com/agent-sandbox/agent-sandbox/pkg/handler"
	"github.com/agent-sandbox/agent-sandbox/pkg/sandbox"
	v1 "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	rsclient "knative.dev/pkg/client/injection/kube/informers/apps/v1/replicaset"
	"knative.dev/pkg/injection"
)

// initTestConfig initializes config for testing if not already initialized
func initTestConfig() {
	if config.Cfg == nil {
		// Set default config values for testing
		config.Cfg = &config.Config{
			APIVersion:             "v1",
			APIBaseURL:             "/api/v1",
			ServerAddr:             "0.0.0.0:10000",
			SandboxNamespace:       "default",
			SandboxDefaultImage:    "ghcr.io/agent-infra/sandbox:latest",
			SandboxDefaultTemplate: "aio",
		}
		config.Cfg.APIBaseURL = "/api/" + config.Cfg.APIVersion
	}
}

// setupMockK8sClient creates a fake k8s client and sets up the injection context
func setupMockK8sClient(ctx context.Context, namespace string) (context.Context, kubernetes.Interface) {
	// Create fake k8s client with reactors to automatically set ReplicaSet status
	fakeClient := fake.NewSimpleClientset()

	// Add reactor to automatically set ReadyReplicas when ReplicaSet is created
	fakeClient.PrependReactor("create", "replicasets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		createAction := action.(ktesting.CreateAction)
		rs := createAction.GetObject().(*v1.ReplicaSet).DeepCopy()
		// Set status to ready immediately for testing
		replicas := int32(1)
		if rs.Spec.Replicas != nil {
			replicas = *rs.Spec.Replicas
		}
		rs.Status = v1.ReplicaSetStatus{
			Replicas:          replicas,
			ReadyReplicas:     replicas,
			AvailableReplicas: replicas,
		}
		return false, rs, nil
	})

	// Add reactor to update status on Get requests (for polling)
	fakeClient.PrependReactor("get", "replicasets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		getAction := action.(ktesting.GetAction)
		rs, err := fakeClient.AppsV1().ReplicaSets(getAction.GetNamespace()).Get(ctx, getAction.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, nil, err
		}
		// Ensure status is set
		replicas := int32(1)
		if rs.Spec.Replicas != nil {
			replicas = *rs.Spec.Replicas
		}
		if rs.Status.ReadyReplicas == 0 {
			rs.Status = v1.ReplicaSetStatus{
				Replicas:          replicas,
				ReadyReplicas:     replicas,
				AvailableReplicas: replicas,
			}
		}
		return false, rs, nil
	})

	// Set up injection context with fake client
	ctx = injection.WithNamespaceScope(ctx, namespace)
	// Inject fake client using context.WithValue (knative injection pattern)
	ctx = context.WithValue(ctx, kubeclient.Key{}, fakeClient)

	// Create informer factory with fake client
	informerFactory := informers.NewSharedInformerFactory(fakeClient, time.Second*30)
	rsInformer := informerFactory.Apps().V1().ReplicaSets()
	// Inject informer using context.WithValue (knative injection pattern)
	ctx = context.WithValue(ctx, rsclient.Key{}, rsInformer)

	// Start informer in background
	stopCh := make(chan struct{})
	go informerFactory.Start(stopCh)

	// Wait a bit for informer to start
	time.Sleep(50 * time.Millisecond)

	return ctx, fakeClient
}

// createMockReplicaSet creates a mock ReplicaSet for testing
func createMockReplicaSet(name, namespace, sandboxID, user string, sandboxData []byte) *v1.ReplicaSet {
	replicas := int32(1)
	return &v1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				sandbox.IDLabel:   sandboxID,
				sandbox.UserLabel: user,
				"sandbox":         name,
				"owner":           "agent-sandbox",
			},
			Annotations: map[string]string{
				"sandbox-data": string(sandboxData),
			},
		},
		Spec: v1.ReplicaSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"sandbox": name,
				},
			},
			Template: v1core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"sandbox": name,
						"owner":   "agent-sandbox",
					},
				},
				Spec: v1core.PodSpec{
					Containers: []v1core.Container{
						{
							Name:  "sandbox",
							Image: "ghcr.io/agent-infra/sandbox:latest",
							StartupProbe: &v1core.Probe{
								ProbeHandler: v1core.ProbeHandler{
									TCPSocket: &v1core.TCPSocketAction{
										Port: intstr.FromInt(8080),
									},
								},
								FailureThreshold:    600,
								PeriodSeconds:       1,
								SuccessThreshold:    1,
								TimeoutSeconds:      3,
								InitialDelaySeconds: 0,
							},
						},
					},
				},
			},
		},
		Status: v1.ReplicaSetStatus{
			Replicas:          replicas,
			ReadyReplicas:     replicas,
			AvailableReplicas: replicas,
		},
	}
}

func TestCreateSandbox(t *testing.T) {
	initTestConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup mock k8s client
	namespace := config.Cfg.SandboxNamespace
	ctx, fakeClient := setupMockK8sClient(ctx, namespace)

	// Create test sandbox request
	testSandbox := &sandbox.Sandbox{
		SandboxBase: sandbox.SandboxBase{
			Name:     "test-sandbox",
			Template: "aio",
		},
		CPU:         "100m",
		Memory:      "128Mi",
		CPULimit:    "1000m",
		MemoryLimit: "1024Mi",
		Timeout:     10,
		IdleTimeout: 10,
		Port:        8080,
	}

	body, err := json.Marshal(testSandbox)
	if err != nil {
		t.Fatalf("Failed to marshal sandbox: %v", err)
	}

	// Create HTTP request
	req := httptest.NewRequest("POST", fmt.Sprintf("%s/sandbox", config.Cfg.APIBaseURL), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Create handler and get the mux
	apiServer := handler.New(ctx)

	// Execute request using the server's handler
	apiServer.Handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Wait a bit for the ReplicaSet to be created (since Create is async in real k8s)
	time.Sleep(100 * time.Millisecond)

	// Verify ReplicaSet was created
	rs, err := fakeClient.AppsV1().ReplicaSets(namespace).Get(ctx, testSandbox.Name, metav1.GetOptions{})
	if err != nil {
		t.Errorf("Failed to get created ReplicaSet: %v", err)
		return
	}

	if rs == nil {
		t.Error("ReplicaSet was not created")
		return
	}

	if rs.Name != testSandbox.Name {
		t.Errorf("Expected ReplicaSet name %s, got %s", testSandbox.Name, rs.Name)
	}

	// Verify sandbox data in annotations
	sandboxData, ok := rs.Annotations["sandbox-data"]
	if !ok {
		t.Error("sandbox-data annotation not found")
	}

	var createdSandbox sandbox.Sandbox
	if err := json.Unmarshal([]byte(sandboxData), &createdSandbox); err != nil {
		t.Errorf("Failed to unmarshal sandbox data: %v", err)
	}

	if createdSandbox.Name != testSandbox.Name {
		t.Errorf("Expected sandbox name %s, got %s", testSandbox.Name, createdSandbox.Name)
	}
}

func TestListSandbox(t *testing.T) {
	initTestConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup mock k8s client
	namespace := config.Cfg.SandboxNamespace
	ctx, fakeClient := setupMockK8sClient(ctx, namespace)

	// Create test sandboxes
	user := "default-user"
	sandbox1 := sandbox.GetDefaultSandbox(user)
	sandbox1.Name = "test-sandbox-1"
	sandbox1.Template = "aio"

	sandbox2 := sandbox.GetDefaultSandbox(user)
	sandbox2.Name = "test-sandbox-2"
	sandbox2.Template = "aio"

	// Create ReplicaSets in fake client
	sandbox1Data, _ := json.Marshal(sandbox1)
	rs1 := createMockReplicaSet(sandbox1.Name, namespace, sandbox1.ID, user, sandbox1Data)
	rs1.Status.ReadyReplicas = 1
	rs1.Status.Replicas = 1

	sandbox2Data, _ := json.Marshal(sandbox2)
	rs2 := createMockReplicaSet(sandbox2.Name, namespace, sandbox2.ID, user, sandbox2Data)
	rs2.Status.ReadyReplicas = 1
	rs2.Status.Replicas = 1

	_, err := fakeClient.AppsV1().ReplicaSets(namespace).Create(ctx, rs1, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create ReplicaSet 1: %v", err)
	}

	_, err = fakeClient.AppsV1().ReplicaSets(namespace).Create(ctx, rs2, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create ReplicaSet 2: %v", err)
	}

	// Wait for informer to sync
	time.Sleep(100 * time.Millisecond)

	// Create HTTP request
	req := httptest.NewRequest("GET", fmt.Sprintf("%s/sandbox", config.Cfg.APIBaseURL), nil)
	w := httptest.NewRecorder()

	// Create handler and get the mux
	apiServer := handler.New(ctx)

	// Execute request using the server's handler
	apiServer.Handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Parse response
	var response struct {
		Code string             `json:"code"`
		Data []*sandbox.Sandbox `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Code != "0" {
		t.Errorf("Expected code '0', got '%s'", response.Code)
	}

	if len(response.Data) < 2 {
		t.Errorf("Expected at least 2 sandboxes, got %d", len(response.Data))
	}
}

func TestGetSandbox(t *testing.T) {
	initTestConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup mock k8s client
	namespace := config.Cfg.SandboxNamespace
	ctx, fakeClient := setupMockK8sClient(ctx, namespace)

	// Create test sandbox
	user := "default-user"
	testSandbox := sandbox.GetDefaultSandbox(user)
	testSandbox.Name = "test-sandbox-get"
	testSandbox.Template = "aio"

	// Create ReplicaSet in fake client
	sandboxData, _ := json.Marshal(testSandbox)
	rs := createMockReplicaSet(testSandbox.Name, namespace, testSandbox.ID, user, sandboxData)
	rs.Status.ReadyReplicas = 1
	rs.Status.Replicas = 1

	_, err := fakeClient.AppsV1().ReplicaSets(namespace).Create(ctx, rs, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create ReplicaSet: %v", err)
	}

	// Create HTTP request
	req := httptest.NewRequest("GET", fmt.Sprintf("%s/sandbox/%s", config.Cfg.APIBaseURL, testSandbox.Name), nil)
	w := httptest.NewRecorder()

	// Create handler and get the mux
	apiServer := handler.New(ctx)

	// Execute request using the server's handler
	apiServer.Handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Parse response
	var response struct {
		Code string           `json:"code"`
		Data *sandbox.Sandbox `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Code != "0" {
		t.Errorf("Expected code '0', got '%s'", response.Code)
	}

	if response.Data == nil {
		t.Error("Expected sandbox data, got nil")
		return
	}

	if response.Data.Name != testSandbox.Name {
		t.Errorf("Expected sandbox name %s, got %s", testSandbox.Name, response.Data.Name)
	}
}

func TestDeleteSandbox(t *testing.T) {
	initTestConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup mock k8s client
	namespace := config.Cfg.SandboxNamespace
	ctx, fakeClient := setupMockK8sClient(ctx, namespace)

	// Create test sandbox
	user := "default-user"
	testSandbox := sandbox.GetDefaultSandbox(user)
	testSandbox.Name = "test-sandbox-delete"
	testSandbox.Template = "aio"

	// Create ReplicaSet in fake client
	sandboxData, _ := json.Marshal(testSandbox)
	rs := createMockReplicaSet(testSandbox.Name, namespace, testSandbox.ID, user, sandboxData)
	rs.Status.ReadyReplicas = 1
	rs.Status.Replicas = 1

	_, err := fakeClient.AppsV1().ReplicaSets(namespace).Create(ctx, rs, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Failed to create ReplicaSet: %v", err)
	}

	// Create HTTP request
	req := httptest.NewRequest("DELETE", fmt.Sprintf("%s/sandbox/%s", config.Cfg.APIBaseURL, testSandbox.Name), nil)
	w := httptest.NewRecorder()

	// Create handler and get the mux
	apiServer := handler.New(ctx)

	// Execute request using the server's handler
	apiServer.Handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify ReplicaSet was deleted
	_, err = fakeClient.AppsV1().ReplicaSets(namespace).Get(ctx, testSandbox.Name, metav1.GetOptions{})
	if err == nil {
		t.Error("Expected ReplicaSet to be deleted, but it still exists")
	}
}

func TestGetSandboxNotFound(t *testing.T) {
	initTestConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup mock k8s client
	namespace := config.Cfg.SandboxNamespace
	ctx, _ = setupMockK8sClient(ctx, namespace)

	// Create HTTP request for non-existent sandbox
	req := httptest.NewRequest("GET", fmt.Sprintf("%s/sandbox/non-existent", config.Cfg.APIBaseURL), nil)
	w := httptest.NewRecorder()

	// Create handler and get the mux
	apiServer := handler.New(ctx)

	// Execute request using the server's handler
	apiServer.Handler.ServeHTTP(w, req)

	// Check response - should return error
	if w.Code == http.StatusOK {
		var response struct {
			Code  string `json:"code"`
			Error string `json:"error"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err == nil {
			if response.Code == "0" {
				t.Error("Expected error response for non-existent sandbox")
			}
		}
	}
}
