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

package e2b

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/agent-sandbox/agent-sandbox/pkg/api/e2b/api"
	"github.com/agent-sandbox/agent-sandbox/pkg/sandbox"
	"github.com/agent-sandbox/agent-sandbox/pkg/utils"
	"k8s.io/klog/v2"
)

func (a *Handler) GetSandbox(w http.ResponseWriter, r *http.Request) {
	sandboxID := r.PathValue("sandboxID")
	if sandboxID == "" {
		sendAPIError(w, http.StatusBadRequest, "sandboxID is required")
		return
	}

	klog.V(2).Infof("Get sandbox sandboxID=%s", sandboxID)

	sbxName := GenSandboxName(sandboxID)

	sb, err := a.controller.Get(sbxName)
	if err != nil {
		klog.Errorf("Get sandbox %s error: %v", sandboxID, err)
		sendAPIError(w, http.StatusNotFound, fmt.Sprintf("sandbox %s not found", sandboxID))
		return
	}

	apiSbx := a.convertToE2BSandbox(sb)
	sendAPIOK(w, http.StatusOK, apiSbx)
	return
}

func (a *Handler) ListSandboxes(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == "" {
		sendAPIError(w, http.StatusBadRequest, "User not found, api key may be invalid")
		return
	}

	sbs, err := a.controller.List(user)
	if err != nil {
		sendAPIError(w, http.StatusInternalServerError, fmt.Sprintf("failed to list sandboxes: %v", err))
		return
	}

	var apiSandboxes = []*api.Sandbox{}
	for _, sb := range sbs {
		apiSbx := a.convertToE2BSandbox(sb)
		apiSandboxes = append(apiSandboxes, apiSbx)
	}

	sendAPIOK(w, http.StatusOK, apiSandboxes)
	return
}

func (a *Handler) DeleteSandbox(w http.ResponseWriter, r *http.Request) {
	sandboxID := r.PathValue("sandboxID")
	if sandboxID == "" {
		sendAPIError(w, http.StatusBadRequest, "sandboxID is required")
		return
	}

	klog.V(2).Infof("Delete sandbox sandboxID=%s", sandboxID)

	sbxName := GenSandboxName(sandboxID)

	err := a.controller.Delete(sbxName)
	if err != nil {
		klog.ErrorS(err, "error deleting sandbox", "sandboxID", sandboxID)
		sendAPIError(w, http.StatusBadRequest, fmt.Sprintf("failed to delete sandbox %s: %v", sandboxID, err))
		return
	}

	sendAPIOK(w, http.StatusOK, nil)
	return
}

func (a *Handler) PostSandboxes(w http.ResponseWriter, r *http.Request) {
	var request *api.NewSandbox

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendAPIError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode request body: %v", err))
		return
	}

	sbx, err := a.CreateSandbox(r.Context(), request)
	if err != nil {
		sendAPIError(w, err.Code, fmt.Sprintf("failed to create sandbox: %v", err.ClientMsg))
		return
	}

	sendAPIOK(w, http.StatusCreated, sbx)
	return
}

func (a *Handler) CreateSandbox(ctx context.Context, newSandbox *api.NewSandbox) (*api.Sandbox, *APIError) {
	user := GetUserFromContext(ctx)
	if user == "" {
		return nil, &APIError{
			ClientMsg: "User not found, api key may be invalid",
			Code:      http.StatusBadRequest,
		}
	}

	var sb = sandbox.GetDefaultSandbox(user)
	accessToken := sb.ID

	//code-interpreter-v1, remove  rightmost  version part of TemplateID in default mode
	tplID := strings.Split(newSandbox.TemplateID, "-v")[0]

	sb.Name = GenSandboxName(sb.ID)
	sb.Template = tplID
	sb.Metadata = newSandbox.Metadata
	sb.EnvVars = newSandbox.EnvVars
	sb.Timeout = newSandbox.Timeout

	klog.Infof("Creating sandbox orgin newSandbox is %+v", newSandbox)

	if err := sb.Make(); err != nil {
		klog.ErrorS(err, "error create sandbox", "sandbox", sb)
		return nil, &APIError{
			Err:       err,
			ClientMsg: "error creating sandbox, error " + err.Error(),
			Code:      http.StatusBadRequest,
		}
	}

	klog.Infof("Creating sandbox with values%+v", sb)

	annotations := sb.GetAnnotations()
	annotations["e2b.envd-access-token"] = accessToken
	sb.SetAnnotations(annotations)

	sbCreated, err := a.controller.Create(sb)
	if err != nil {
		klog.ErrorS(err, "error creating sandbox", "sandbox", sb)

		return nil, &APIError{
			Err:       err,
			ClientMsg: "error creating sandbox, error " + err.Error(),
			Code:      http.StatusBadRequest,
		}
	}

	apiSbx := a.convertToE2BSandbox(sbCreated)

	return apiSbx, nil
}

func (a *Handler) ConnectSandbox(w http.ResponseWriter, r *http.Request) {
	sandboxID := r.PathValue("sandboxID")
	if sandboxID == "" {
		sendAPIError(w, http.StatusBadRequest, "sandboxID is required")
		return
	}

	sbxName := GenSandboxName(sandboxID)

	sb, err := a.controller.Get(sbxName)
	if err != nil {
		klog.Errorf("Get sandbox %s error: %v", sandboxID, err)
		sendAPIError(w, http.StatusNotFound, fmt.Sprintf("sandbox %s not found", sandboxID))
		return
	}

	//TODO active the sandbox if it's paused or stopped

	apiSbx := a.convertToE2BSandbox(sb)

	sendAPIOK(w, http.StatusCreated, apiSbx)
	return
}

func (a *Handler) SandboxRouterOfPath() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		klog.Info("Entering SandboxRouterOfPath", " url=", r.URL.Path, "method=", r.Method, "query=", r.URL.RawQuery)

		sandboxID := r.PathValue("sandboxID")
		if sandboxID == "" {
			sendAPIError(w, http.StatusBadRequest, "sandboxID is required")
			return
		}

		port := r.PathValue("port")
		if port == "" {
			sendAPIError(w, http.StatusBadRequest, "port is required")
			return
		}

		sbName := GenSandboxName(sandboxID)

		// rewrite url to /sandbox/{name}/{port}/...
		prefixToStrip := fmt.Sprintf("/sandboxes/router/%s/%s", sandboxID, port)
		postPath := strings.TrimPrefix(r.URL.Path, prefixToStrip)

		pxyURL := fmt.Sprintf("/sandbox/%s%s", sbName, postPath)
		r.URL.Path = pxyURL
		//add port query param to url
		if r.URL.RawQuery == "" {
			r.URL.RawQuery = "port=" + port
		} else {
			r.URL.RawQuery = r.URL.RawQuery + "&port=" + port
		}

		klog.Info("ExecuteSandbox proxying... new url=", r.URL)
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}
}

func (a *Handler) SandboxRouterNative() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		klog.Info("Entering SandboxRouterNative", " url=", r.URL, "method=", r.Method, "query=", r.URL.RawQuery)

		// match host and request url is {port}-{sandbox_id}.{domain}/...
		reqHost := r.Host
		//"49999-e94466d4e94466d4.example.com"
		patten := `^[0-9]{4,5}-[a-f0-9]{8,48}\.[a-zA-Z0-9.-]+`
		matched, _ := regexp.MatchString(patten, reqHost)
		if !matched {
			// otherwise just return 404
			klog.Info("no matched sandbox proxy pattern, url=", reqHost)
			sendAPIError(w, http.StatusBadRequest, "invalid sandbox request pattern 1")
			return
		}

		// match port from reqURL
		portSvr := regexp.MustCompile(`^([0-9]{4,5})-([a-f0-9]+)\.`)
		submatches := portSvr.FindStringSubmatch(reqHost)
		if len(submatches) < 3 {
			klog.Info("invalid matched sandbox request pattern, url=", reqHost)
			sendAPIError(w, http.StatusBadRequest, "invalid sandbox request pattern 2")
			return
		}
		port := submatches[1]
		sandboxID := submatches[2]

		sbName := GenSandboxName(sandboxID)

		// rewrite url to /sandbox/{name}/{port}/...
		pxyURL := fmt.Sprintf("/sandbox/%s%s", sbName, r.URL.Path)
		r.URL.Path = pxyURL
		if r.URL.RawQuery == "" {
			r.URL.RawQuery = "port=" + port
		} else {
			r.URL.RawQuery = r.URL.RawQuery + "&port=" + port
		}

		klog.Info("matched sandbox proxy pattern, proxying... new url=", r.URL)
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}
}

// convertToE2BSandbox converts an internal sandbox.Sandbox to api.Sandbox format
func (a *Handler) convertToE2BSandbox(sb *sandbox.Sandbox) *api.Sandbox {
	apiSbx := GetDefaultE2BSandbox()
	apiSbx.SandboxID = sb.ID
	apiSbx.TemplateID = sb.Template

	annotations := sb.GetAnnotations()
	if accessToken, ok := annotations["e2b.envd-access-token"]; ok {
		apiSbx.EnvdAccessToken = accessToken
		apiSbx.TrafficAccessToken = accessToken
	}

	if sb.Metadata != nil {
		apiSbx.Metadata = sb.Metadata
	}

	rs := utils.CalculateResourceToQuantity(sb)
	apiSbx.CpuCount = rs.CPUMilli
	apiSbx.MemoryMB = rs.MemoryMB
	apiSbx.DiskSizeMB = rs.DiskSizeMB

	apiSbx.StartedAt = sb.GetCreationTimestamp().Time

	apiSbx.State = api.Running
	if sb.Status != "" {
		apiSbx.State = api.SandboxState(sb.Status)
	}

	return apiSbx
}

func GenSandboxName(id string) string {
	name := fmt.Sprintf("sandbox-e2b-%s", id)

	return name
}
