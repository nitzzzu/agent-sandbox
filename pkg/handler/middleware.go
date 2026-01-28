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
	"net/http"
	"time"

	"k8s.io/klog/v2"
)

// LoggingMiddleware logs the request details and execution time
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//skip healthz logging
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Code to execute BEFORE the actual handler
		klog.V(1).Infof("Started request: %s %s%s, request headers %v", r.Method, r.Host, r.URL, r.Header)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		klog.V(1).Infof("Finished request: %s %s%s, Time taken: %v, response header: %v", r.Method, r.Host, r.URL, time.Since(start), w.Header())
	})
}
