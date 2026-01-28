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

package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"k8s.io/klog/v2"
)

type Template struct {
	Name        string `json:"name" required:"false"`
	Image       string `json:"image" required:"false"`
	Port        int    `json:"port" required:"false"`
	Description string `json:"description" required:"false"`
}

var Cfg *Config
var Templates *[]*Template

type Config struct {
	APIVersion string `split_words:"true" default:"v1" required:"false"`
	APIBaseURL string `split_words:"true" default:"" required:"false"`
	ServerAddr string `split_words:"true" default:"0.0.0.0:10000" required:"false"`

	// witch Kubernetes namespace to create sandboxes Replicaset&Pod in
	SandboxNamespace string `split_words:"true" default:"default" required:"false"`

	SandboxTemplateFile string `split_words:"true" default:"" required:"false"`

	SandboxTemplatesConfigFile string `split_words:"true" default:"templates.json" required:"false"`
	SandboxDefaultImage        string `split_words:"true" default:"ghcr.io/agent-infra/sandbox:latest" required:"false"`
	SandboxDefaultTemplate     string `split_words:"true" default:"aio" required:"false"`
}

func init() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		klog.Fatal("Failed to process config: ", err)
	}

	cfg.APIBaseURL = "/api/" + cfg.APIVersion
	Cfg = &cfg

	LoadTemplates()
}

func LoadTemplates() {
	//load templates config, read file from cfg.SandboxTemplatesConfigFile by os.ReadFile
	envFile := Cfg.SandboxTemplatesConfigFile
	klog.Infof("Loading Template config from file %s", envFile)

	templates, err := os.ReadFile(envFile)
	if err != nil {
		klog.Fatalf("Failed to read Template config file %s error: %v", envFile, err)
	}

	klog.Infof("Loaded Template config from file %s, content is %s", envFile, string(templates))

	var tpls []*Template
	err = json.Unmarshal(templates, &tpls)
	if err != nil {
		klog.Fatalf("Failed to unmarshal Template config file %s error: %v", envFile, err)
	}

	//check envs not empty
	if len(tpls) == 0 {
		klog.Fatalf("No Templates  found in config file %s", envFile)
	}

	//varify each env has name  image and description
	for _, env := range tpls {
		if env.Name == "" || env.Image == "" || env.Description == "" {
			klog.Fatalf("Invalid Template config in file %s: %+v, name image and desc must not dempty", envFile, env)
		}
	}

	Templates = &tpls

	//log loaded envs
	for _, env := range *Templates {
		klog.Infof("Loaded Template object: %+v", env)
	}
}

func GetTemplateByName(name string) (*Template, error) {
	for _, env := range *Templates {
		if env.Name == name {
			return env, nil
		}
	}
	klog.Errorf("Template %s not found", name)
	return nil, fmt.Errorf("Template  %s not found", name)
}

// GetTemplatesForMCPTools return json string, but exclude image field
func GetTemplatesForMCPTools() string {
	type TplForTool struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var tpls []TplForTool
	for _, env := range *Templates {
		tpls = append(tpls, TplForTool{
			Name:        env.Name,
			Description: env.Description,
		})
	}

	tplsJson, err := json.MarshalIndent(tpls, "", "  ")
	if err != nil {
		klog.Errorf("Failed to marshal Templates for MCP tools: %v", err)
		return ""
	}

	return string(tplsJson)
}
