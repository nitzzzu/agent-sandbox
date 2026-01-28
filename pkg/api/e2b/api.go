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
	"fmt"
	"net/http"
	"os"

	"github.com/agent-sandbox/agent-sandbox/pkg/activator"
	"github.com/agent-sandbox/agent-sandbox/pkg/api/e2b/api"
	"github.com/agent-sandbox/agent-sandbox/pkg/sandbox"

	"context"
)

const (
	BaseURL = "/e2b/v1"
)

type Handler struct {
	rootCtx    context.Context
	controller *sandbox.Controller
	activator  *activator.Activator
	domain     string
}

func NewHandler(rootCtx context.Context, c *sandbox.Controller, a *activator.Activator) *Handler {
	domain := "localhost"
	if domainEnv := os.Getenv("E2B_DOMAIN"); domainEnv != "" {
		domain = domainEnv
	}

	return &Handler{
		rootCtx:    rootCtx,
		controller: c,
		activator:  a,
		domain:     domain,
	}
}

func GetDefaultE2BSandbox() *api.Sandbox {
	return &api.Sandbox{
		EnvdVersion:     "0.1.1",
		EnvdAccessToken: "envd-access-token-x",
		ClientID:        "client-id-x",
		TemplateID:      "code-interpreter-v1",
	}
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func (a *Handler) RegisterHandlersWithOptions(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("POST %s/sandboxes", BaseURL), a.withCheckApiKey(a.PostSandboxes))
	mux.HandleFunc(fmt.Sprintf("GET %s/sandboxes/{sandboxID}", BaseURL), a.withCheckApiKey(a.GetSandbox))
	mux.HandleFunc(fmt.Sprintf("GET %s/v2/sandboxes", BaseURL), a.withCheckApiKey(a.ListSandboxes))

	//E2B does not support delete sandbox now, delete API is reserved for call with http api
	mux.HandleFunc(fmt.Sprintf("DELETE %s/sandboxes/{sandboxID}", BaseURL), a.withCheckApiKey(a.DeleteSandbox))

	mux.HandleFunc(fmt.Sprintf("POST %s/sandboxes/{sandboxID}/connect", BaseURL), a.withCheckApiKey(a.ConnectSandbox))

}

// withCheckApiKey wraps an http.HandlerFunc with CheckApiKey middleware
func (a *Handler) withCheckApiKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, apiErr := a.CheckApiKey(r.Context(), r)
		if apiErr != nil {
			sendAPIError(w, apiErr.Code, apiErr.ClientMsg)
			return
		}
		// Update request with the new context
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func (a *Handler) CheckApiKey(ctx context.Context, r *http.Request) (context.Context, *APIError) {
	// e.g. e2b-aef134ef-7aa1-945e-9399-7df9a4ad0c3f
	apiKey := r.Header.Get("X-Api-Key")

	//TODO check api key validity from database or config

	if apiKey == "" {
		return ctx, &APIError{
			Code:      http.StatusUnauthorized,
			ClientMsg: "Missing API Key",
		}
	}

	return context.WithValue(ctx, "user", apiKey), nil
}

func GetUserFromContext(ctx context.Context) string {
	value := ctx.Value("user")
	user, ok := value.(string)
	if !ok {
		return ""
	}
	return user
}
