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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kelseyhightower/envconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type TemplatePool struct {
	Size       int    `json:"size"`
	ProbePort  int    `json:"probePort"`
	WarmupCmd  string `json:"warmupCmd"`
	StartupCmd string `json:"startupCmd"`
}

type Template struct {
	Name           string       `json:"name" required:"false"`
	Pattern        string       `json:"pattern" required:"false"`
	Image          string       `json:"image" required:"false"`
	Port           int          `json:"port" required:"false"`
	Type           string       `json:"type" required:"false" description:"dynamic or static, default is static, dynamic means template is dynamic by regexp"`
	NoStartupProbe bool         `json:"noStartupProbe" required:"false"`
	Pool           TemplatePool `json:"pool" required:"false"`
	Description    string       `json:"description" required:"false"`
}

var Cfg *Config
var Templates []*Template
var SandboxDeployTemplate string

type Config struct {
	KubeClient kubernetes.Interface `ignored:"true"`

	APIVersion string `split_words:"true" default:"v1" required:"false"`
	APIBaseURL string `split_words:"true" default:"" required:"false"`
	ServerAddr string `split_words:"true" default:"0.0.0.0:10000" required:"false"`

	// witch Kubernetes namespace to create sandboxes Replicaset&Pod in
	SandboxNamespace string `split_words:"true" default:"default" required:"false"`

	SandboxTemplateFile string `split_words:"true" default:"config/sandbox.yaml" required:"false"`

	SandboxTemplatesConfigFile string `split_words:"true" default:"config/templates.json" required:"false"`
	SandboxDefaultImage        string `split_words:"true" default:"ghcr.io/agent-infra/sandbox:latest" required:"false"`
	SandboxDefaultTemplate     string `split_words:"true" default:"aio" required:"false"`
}

func InitConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		klog.Fatal("Failed to process config: ", err)
	}

	cfg.APIBaseURL = "/api/" + cfg.APIVersion
	Cfg = &cfg

	return Cfg
}

func (c *Config) LoadSandboxRSTemplate() {
	var err error
	var val []byte
	if val, err = os.ReadFile(c.SandboxTemplateFile); err != nil {
		panic(err)
	}
	SandboxDeployTemplate = string(val)
}

func (c *Config) CheckConfigmap() {
	cmcontent, err := c.ReadTemplatesFromCM()
	if cmcontent != "" {
		klog.Info("templates config already exists in configmap")
		return
	}

	klog.Info("templates config is empty, will load from local file system")

	//load templates config, read file from cfg.SandboxTemplatesConfigFile by os.ReadFile
	fileName := c.SandboxTemplatesConfigFile
	content, err := os.ReadFile(fileName)
	if err != nil {
		klog.Fatalf("Failed to read Template config file %s error: %v", fileName, err)
	}

	templatesContent := string(content)

	klog.Infof("Loaded Template config from file %s, content is %s", fileName, templatesContent)

	// store templates config to configmap
	err = c.SaveTemplatesToCM(templatesContent)
	if err != nil {
		klog.Fatalf("Failed to save Templates config to configmap, error: %v, config content %v", err, templatesContent)
	} else {
		klog.Info("Templates config saved to configmap successfully")
	}
}

// LoadTemplates load templates config from:
func (c *Config) LoadTemplates() {
	templatesContent := ""

	// load config from configmap
	content, err := c.ReadTemplatesFromCM()
	if content == "" {
		klog.Fatalf("Failed to read Templates config from configmap, content is empty, error: %v", err)
	}
	klog.Info("Loaded Templates config from configmap: ", content)
	templatesContent = content

	var tpls []*Template
	err = json.Unmarshal([]byte(templatesContent), &tpls)
	if err != nil {
		klog.Fatalf("Failed to unmarshal Template config templatesContent %s error: %v", templatesContent, err)
	}

	//check envs not empty
	if len(tpls) == 0 {
		klog.Fatalf("No Templates  found in config content %s", templatesContent)
	}

	//varify each env has name  image and description
	for _, env := range tpls {
		if env.Name == "" || env.Image == "" || env.Description == "" {
			klog.Fatalf("Invalid Template config in templatesContent %s: %+v, name image and desc must not dempty", templatesContent, env)
		}
	}

	Templates = tpls

	//log loaded envs
	for _, env := range Templates {
		klog.Infof("Loaded Template object: %+v", env)
	}
}

func (c *Config) SaveTemplatesToCM(templatesContent string) error {
	var err error
	// store templates config to configmap
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TemplatesConfigMapName,
			Namespace: c.SandboxNamespace,
		},
		Data: map[string]string{
			TemplatesConfigMapKey: templatesContent,
		},
	}
	_, err = c.KubeClient.CoreV1().ConfigMaps(c.SandboxNamespace).Create(context.TODO(), cm, metav1.CreateOptions{})

	return err
}

func (c *Config) ReadTemplatesFromCM() (content string, err error) {
	templatesContent := ""

	existCm, err := c.KubeClient.CoreV1().ConfigMaps(c.SandboxNamespace).Get(context.TODO(), TemplatesConfigMapName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		klog.Info("templates configmap not found, return empty content")
		return templatesContent, nil
	} else if err != nil {
		klog.Errorf("Failed to get ConfigMap for Templates config: %v", err)
		return templatesContent, err
	}

	return existCm.Data[TemplatesConfigMapKey], nil

}

func GetTemplateByName(name string) (*Template, error) {
	for _, t := range Templates {
		if t.Name == name {
			return t, nil
		}
	}

	for _, t := range Templates {
		if t.Type == "dynamic" {
			image := t.Image
			//match by regexp
			re := regexp.MustCompile(t.Pattern)
			match := re.FindStringSubmatch(name)
			if len(match) == 0 {
				continue
			}

			if len(match) > 0 {
				versionIndex := re.SubexpIndex("version")
				nameIndex := re.SubexpIndex("name")
				if nameIndex == -1 || versionIndex == -1 {
					continue
				}

				tversion := match[versionIndex]
				tname := match[nameIndex]
				image = strings.ReplaceAll(image, "<name>", tname)
				image = strings.ReplaceAll(image, "<version>", tversion)
			}

			dynT := &Template{
				Name:           t.Name,
				Image:          image,
				Port:           t.Port,
				Pattern:        t.Pattern,
				Pool:           t.Pool,
				Type:           t.Type,
				NoStartupProbe: t.NoStartupProbe,
				Description:    t.Description,
			}
			return dynT, nil
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
	for _, env := range Templates {
		if env.Type == "dynamic" {
			continue
		}
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
