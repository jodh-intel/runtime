// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package virtcontainers

import (
	"syscall"

	"github.com/Sirupsen/logrus"
)

// VC is the Virtcontainers interface
type VC interface {
	SetLogger(logger logrus.FieldLogger)

	CreatePod(podConfig PodConfig) (VCPod, error)
	DeletePod(podID string) (VCPod, error)
	StartPod(podID string) (VCPod, error)
	StopPod(podID string) (VCPod, error)
	RunPod(podConfig PodConfig) (VCPod, error)
	ListPod() ([]PodStatus, error)
	StatusPod(podID string) (PodStatus, error)
	PausePod(podID string) (VCPod, error)
	ResumePod(podID string) (VCPod, error)

	CreateContainer(podID string, containerConfig ContainerConfig) (VCPod, VCContainer, error)
	DeleteContainer(podID, containerID string) (VCContainer, error)
	StartContainer(podID, containerID string) (VCContainer, error)
	StopContainer(podID, containerID string) (VCContainer, error)
	EnterContainer(podID, containerID string, cmd Cmd) (VCPod, VCContainer, *Process, error)
	StatusContainer(podID, containerID string) (ContainerStatus, error)
	KillContainer(podID, containerID string, signal syscall.Signal, all bool) error
}

// VCPod is the Pod interface
// (required since virtcontainers.Pod only contains private fields)
type VCPod interface {
	ID() string
	Annotations(key string) (string, error)
	SetAnnotations(annotations map[string]string) error
	GetAnnotations() map[string]string
	URL() string
	GetAllContainers() []VCContainer
	GetContainer(containerID string) VCContainer
}

// VCContainer is the Container interface
// (required since virtcontainers.Container only contains private fields)
type VCContainer interface {
	ID() string
	Pod() VCPod
	Process() Process
	GetToken() string
	GetPid() int
	SetPid(pid int) error
	URL() string
	GetAnnotations() map[string]string
}
