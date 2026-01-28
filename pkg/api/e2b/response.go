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

import "net/http"
import "encoding/json"

type APIError struct {
	// Err Internal error, for logging
	Err error

	// ClientMsg Message to be shown to client, for user
	ClientMsg string

	Code int
}

type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error
	Message string `json:"message"`
}

func sendAPIError(w http.ResponseWriter, code int, message string) {
	apiErr := Error{
		Code:    int32(code),
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiErr)
}

func sendAPIOK(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
