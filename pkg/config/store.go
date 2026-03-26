package config

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

const TemplatesConfigMapName = "agent-sandbox"
const TemplatesConfigMapKey = "config-templates"

// SandboxTemplateConfigMapKey sandbox template is k8s Resource Definition(ReplicasSet) for sandbox
const SandboxTemplateConfigMapKey = "config-sandbox-template"

func WatchConfigMap() func(configMap *corev1.ConfigMap) {
	var lastTemplatesContent string
	var lastSandboxTemplateContent string

	return func(configMap *corev1.ConfigMap) {
		templatesContent := configMap.Data[TemplatesConfigMapKey]
		if lastTemplatesContent == "" || templatesContent != lastTemplatesContent {
			klog.Info("watching ConfigMap changed, templates content updated")
			Cfg.ShouldLoadTemplates()
			lastTemplatesContent = templatesContent
		}

		sandboxTemplateContent := configMap.Data[SandboxTemplateConfigMapKey]
		if lastSandboxTemplateContent == "" || sandboxTemplateContent != lastSandboxTemplateContent {
			klog.Info("watching ConfigMap changed, sandbox template content updated")
			Cfg.ShouldLoadSandboxTemplate()
			lastSandboxTemplateContent = sandboxTemplateContent
		}
	}
}
