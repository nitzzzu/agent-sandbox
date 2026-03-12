package config

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

const TemplatesConfigMapName = "agent-sandbox"
const TemplatesConfigMapKey = "config-templates"

func WatchConfigMap() func(configMap *corev1.ConfigMap) {
	return func(configMap *corev1.ConfigMap) {
		content := configMap.Data[TemplatesConfigMapKey]
		if content == "" {
			return
		}
		klog.Info("watching ConfigMap changed, content ", content)
		Cfg.LoadTemplates()
	}
}
