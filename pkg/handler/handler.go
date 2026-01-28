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

package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/agent-sandbox/agent-sandbox/pkg/activator"
	e2bapi "github.com/agent-sandbox/agent-sandbox/pkg/api/e2b"
	"github.com/agent-sandbox/agent-sandbox/pkg/config"
	"github.com/agent-sandbox/agent-sandbox/pkg/router"
	"github.com/agent-sandbox/agent-sandbox/pkg/sandbox"
	"k8s.io/klog/v2"
)

type ApiHttpHandler struct {
	mux        *http.ServeMux
	rootCtx    context.Context
	activator  *activator.Activator
	controller *sandbox.Controller
}

func New(rootCtx context.Context, a *activator.Activator, c *sandbox.Controller) *http.Server {
	mux := http.DefaultServeMux

	ah := &ApiHttpHandler{
		rootCtx:    rootCtx,
		mux:        mux,
		activator:  a,
		controller: c,
	}
	ah.regHandlers()

	// Wrap mux with global logging middleware
	loggedMux := LoggingMiddleware(mux)
	//loggedMux := mux

	server := &http.Server{
		Addr:         config.Cfg.ServerAddr,
		Handler:      loggedMux,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 30 * time.Second,
	}

	klog.Info("Api server ", "addr=", config.Cfg.ServerAddr)
	return server
}

func (ahh *ApiHttpHandler) regHandlers() {
	a := ahh.activator
	c := ahh.controller

	// Rest API for Sandbox management
	sbHeader := sandbox.NewHandler(ahh.rootCtx, c, a)
	ahh.mux.HandleFunc(fmt.Sprintf("POST %s/sandbox", config.Cfg.APIBaseURL), func(w http.ResponseWriter, r *http.Request) { wrapperHandler(w, r, sbHeader.CreateSandbox) })
	ahh.mux.HandleFunc(fmt.Sprintf("GET %s/sandbox", config.Cfg.APIBaseURL), func(w http.ResponseWriter, r *http.Request) { wrapperHandler(w, r, sbHeader.ListSandbox) })
	ahh.mux.HandleFunc(fmt.Sprintf("DELETE %s/sandbox/{name}", config.Cfg.APIBaseURL), func(w http.ResponseWriter, r *http.Request) { wrapperHandler(w, r, sbHeader.DelSandbox) })
	ahh.mux.HandleFunc(fmt.Sprintf("GET %s/sandbox/{name}", config.Cfg.APIBaseURL), func(w http.ResponseWriter, r *http.Request) { wrapperHandler(w, r, sbHeader.GetSandbox) })

	// e2b API
	e2bHeader := e2bapi.NewHandler(ahh.rootCtx, c, a)
	e2bHeader.RegisterHandlersWithOptions(ahh.mux)

	// SandboxHandler router, route calls to Sandbox container
	srHandler := router.NewSandboxRouter(ahh.rootCtx, a)
	ahh.mux.HandleFunc("/sandbox/{name}/", srHandler.ServeHTTP)

	ahh.mux.Handle("/mcp", sbHeader.McpSseHandler())

	ahh.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
		return
	})

	// e2b sandbox execute endpoint
	ahh.mux.HandleFunc("/sandboxes/router/{sandboxID}/{port}/", e2bHeader.SandboxRouterOfPath())
	// catch-all handler for any unmatched requests, mainly for E2B sandbox proxy purpose
	ahh.mux.HandleFunc("/", e2bHeader.SandboxRouterNative())
}

func wrapperHandler(w http.ResponseWriter, r *http.Request, f func(*http.Request) (interface{}, error)) {
	result, err := f(r)
	if err != nil {
		Err(w, err.Error())
		return
	}
	Ok(w, result)
	return
}
