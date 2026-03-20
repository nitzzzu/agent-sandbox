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

package utils

import (
	"bytes"
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
)

func ExecCommand(client kubernetes.Interface, config *rest.Config, namespace, podName, containerName string, command []string) (output string, outputErr string, err error) {
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	option := &corev1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    false,
		TTY:       false,
	}
	req.VersionedParams(option, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", "", err
	}
	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		klog.Errorf("Failed to execute command in pod, namespace: %s, pod: %s, container: %s, command: %v, error: %v, stderr: %s",
			namespace, podName, containerName, command, err, stderr.String())
	}

	return stdout.String(), stderr.String(), err
}

func ExecStreamCommand(ctx context.Context, client kubernetes.Interface, config *rest.Config, namespace, podName, containerName string, command []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	option := &corev1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdin:     stdin != nil,
		Stdout:    stdout != nil,
		Stderr:    stderr != nil,
		TTY:       false,
	}
	req.VersionedParams(option, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	return exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    false,
	})
}

type TerminalSize struct {
	Width  uint16
	Height uint16
}

type TerminalSizeQueue struct {
	resizeCh <-chan TerminalSize
}

func NewTerminalSizeQueue(resizeCh <-chan TerminalSize) remotecommand.TerminalSizeQueue {
	return &TerminalSizeQueue{resizeCh: resizeCh}
}

func (q *TerminalSizeQueue) Next() *remotecommand.TerminalSize {
	size, ok := <-q.resizeCh
	if !ok {
		return nil
	}

	return &remotecommand.TerminalSize{
		Width:  size.Width,
		Height: size.Height,
	}
}

func ExecInteractiveCommand(ctx context.Context, client kubernetes.Interface, config *rest.Config, namespace, podName, containerName string, command []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, resizeQueue remotecommand.TerminalSizeQueue) error {
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	option := &corev1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    false,
		TTY:       true,
	}
	req.VersionedParams(option, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	return exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             stdin,
		Stdout:            stdout,
		Stderr:            stderr,
		Tty:               true,
		TerminalSizeQueue: resizeQueue,
	})
}
