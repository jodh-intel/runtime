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

// Description: A mock implementation of virtcontainers that can be used
// for testing.
//
// This implementation calls the function set in the object that
// corresponds to the name of the method (for example, when CreatePod()
// is called, that method will try to call CreatePodFunc). If no
// function is defined for the method, it will return an error in a
// well-known format. Callers can detect this scenario by calling
// IsMockError().

package vcMock

import (
	"fmt"
	"syscall"

	"github.com/Sirupsen/logrus"
	vc "github.com/containers/virtcontainers"
)

// VCMockPod is a fake Pod type used for testing
type VCMockPod struct {
	MockID          string
	MockURL         string
	MockAnnotations map[string]string
	MockContainers  []*VCMockContainer
}

// VCMockContainer is a fake Container type used for testing
type VCMockContainer struct {
	MockID          string
	MockURL         string
	MockToken       string
	MockProcess     vc.Process
	MockPid         int
	MockPod         *VCMockPod
	MockAnnotations map[string]string
}

// VCMock is an implementation of the VC interface used for testing
type VCMock struct {
	SetLoggerFunc       func(logger logrus.FieldLogger)
	CreatePodFunc       func(podConfig vc.PodConfig) (vc.VCPod, error)
	DeletePodFunc       func(podID string) (vc.VCPod, error)
	StartPodFunc        func(podID string) (vc.VCPod, error)
	StopPodFunc         func(podID string) (vc.VCPod, error)
	RunPodFunc          func(podConfig vc.PodConfig) (vc.VCPod, error)
	ListPodFunc         func() ([]vc.PodStatus, error)
	StatusPodFunc       func(podID string) (vc.PodStatus, error)
	PausePodFunc        func(podID string) (vc.VCPod, error)
	ResumePodFunc       func(podID string) (vc.VCPod, error)
	CreateContainerFunc func(podID string, containerConfig vc.ContainerConfig) (vc.VCPod, vc.VCContainer, error)
	DeleteContainerFunc func(podID, containerID string) (vc.VCContainer, error)
	StartContainerFunc  func(podID, containerID string) (vc.VCContainer, error)
	StopContainerFunc   func(podID, containerID string) (vc.VCContainer, error)
	EnterContainerFunc  func(podID, containerID string, cmd vc.Cmd) (vc.VCPod, vc.VCContainer, *vc.Process, error)
	StatusContainerFunc func(podID, containerID string) (vc.ContainerStatus, error)
	KillContainerFunc   func(podID, containerID string, signal syscall.Signal, all bool) error
}

// mockErrorPrefix is a string that all errors returned by the mock
// implementation itself will contain as a prefix.
const mockErrorPrefix = "vcMock forced failure"

func (m *VCMock) SetLogger(logger logrus.FieldLogger) {
	if m.SetLoggerFunc != nil {
		m.SetLoggerFunc(logger)
		return
	}
}

func (m *VCMock) CreatePod(podConfig vc.PodConfig) (vc.VCPod, error) {
	if m.CreatePodFunc != nil {
		return m.CreatePodFunc(podConfig)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podConfig: %v", mockErrorPrefix, getSelf(), m, podConfig)
}

func (m *VCMock) DeletePod(podID string) (vc.VCPod, error) {
	if m.DeletePodFunc != nil {
		return m.DeletePodFunc(podID)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podID: %v", mockErrorPrefix, getSelf(), m, podID)
}

func (m *VCMock) StartPod(podID string) (vc.VCPod, error) {
	if m.StartPodFunc != nil {
		return m.StartPodFunc(podID)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podID: %v", mockErrorPrefix, getSelf(), m, podID)
}

func (m *VCMock) StopPod(podID string) (vc.VCPod, error) {
	if m.StopPodFunc != nil {
		return m.StopPodFunc(podID)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podID: %v", mockErrorPrefix, getSelf(), m, podID)
}

func (m *VCMock) RunPod(podConfig vc.PodConfig) (vc.VCPod, error) {
	if m.RunPodFunc != nil {
		return m.RunPodFunc(podConfig)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podConfig: %v", mockErrorPrefix, getSelf(), m, podConfig)
}

func (m *VCMock) ListPod() ([]vc.PodStatus, error) {
	if m.ListPodFunc != nil {
		return m.ListPodFunc()
	}

	return nil, fmt.Errorf("%s: %s", mockErrorPrefix, getSelf())
}

func (m *VCMock) StatusPod(podID string) (vc.PodStatus, error) {
	if m.StatusPodFunc != nil {
		return m.StatusPodFunc(podID)
	}

	return vc.PodStatus{}, fmt.Errorf("%s: %s (%+v): podID: %v", mockErrorPrefix, getSelf(), m, podID)
}

func (m *VCMock) PausePod(podID string) (vc.VCPod, error) {
	if m.PausePodFunc != nil {
		return m.PausePodFunc(podID)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podID: %v", mockErrorPrefix, getSelf(), m, podID)
}

func (m *VCMock) ResumePod(podID string) (vc.VCPod, error) {
	if m.ResumePodFunc != nil {
		return m.ResumePodFunc(podID)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podID: %v", mockErrorPrefix, getSelf(), m, podID)
}

func (m *VCMock) CreateContainer(podID string, containerConfig vc.ContainerConfig) (vc.VCPod, vc.VCContainer, error) {
	if m.CreateContainerFunc != nil {
		return m.CreateContainerFunc(podID, containerConfig)
	}

	return nil, nil, fmt.Errorf("%s: %s (%+v): podID: %v, containerConfig: %v", mockErrorPrefix, getSelf(), m, podID, containerConfig)
}

func (m *VCMock) DeleteContainer(podID, containerID string) (vc.VCContainer, error) {
	if m.DeleteContainerFunc != nil {
		return m.DeleteContainerFunc(podID, containerID)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podID: %v, containerID: %v", mockErrorPrefix, getSelf(), m, podID, containerID)
}

func (m *VCMock) StartContainer(podID, containerID string) (vc.VCContainer, error) {
	if m.StartContainerFunc != nil {
		return m.StartContainerFunc(podID, containerID)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podID: %v, containerID: %v", mockErrorPrefix, getSelf(), m, podID, containerID)
}

func (m *VCMock) StopContainer(podID, containerID string) (vc.VCContainer, error) {
	if m.StopContainerFunc != nil {
		return m.StopContainerFunc(podID, containerID)
	}

	return nil, fmt.Errorf("%s: %s (%+v): podID: %v, containerID: %v", mockErrorPrefix, getSelf(), m, podID, containerID)
}

func (m *VCMock) EnterContainer(podID, containerID string, cmd vc.Cmd) (vc.VCPod, vc.VCContainer, *vc.Process, error) {
	if m.EnterContainerFunc != nil {
		return m.EnterContainerFunc(podID, containerID, cmd)
	}

	return nil, nil, nil, fmt.Errorf("%s: %s (%+v): podID: %v, containerID: %v, cmd: %v", mockErrorPrefix, getSelf(), m, podID, containerID, cmd)
}

func (m *VCMock) StatusContainer(podID, containerID string) (vc.ContainerStatus, error) {
	if m.StatusContainerFunc != nil {
		return m.StatusContainerFunc(podID, containerID)
	}

	return vc.ContainerStatus{}, fmt.Errorf("%s: %s (%+v): podID: %v, containerID: %v", mockErrorPrefix, getSelf(), m, podID, containerID)
}

func (m *VCMock) KillContainer(podID, containerID string, signal syscall.Signal, all bool) error {
	if m.KillContainerFunc != nil {
		return m.KillContainerFunc(podID, containerID, signal, all)
	}

	return fmt.Errorf("%s: %s (%+v): podID: %v, containerID: %v, signal: %v, all: %v", mockErrorPrefix, getSelf(), m, podID, containerID, signal, all)
}

func (p *VCMockPod) ID() string {
	return p.MockID
}

func (p *VCMockPod) Annotations(key string) (string, error) {
	return p.MockAnnotations[key], nil
}

func (p *VCMockPod) SetAnnotations(annotations map[string]string) error {
	return nil
}

func (p *VCMockPod) GetAnnotations() map[string]string {
	return p.MockAnnotations
}

func (p *VCMockPod) URL() string {
	return p.MockURL
}

func (p *VCMockPod) GetAllContainers() []vc.VCContainer {
	var ifa []vc.VCContainer = make([]vc.VCContainer, len(p.MockContainers))

	for i, v := range p.MockContainers {
		ifa[i] = v
	}

	return ifa
}

func (p *VCMockPod) GetContainer(containerID string) vc.VCContainer {
	for _, c := range p.MockContainers {
		if c.MockID == containerID {
			return c
		}
	}
	return &VCMockContainer{}
}

func (c *VCMockContainer) ID() string {
	return c.MockID
}

func (c *VCMockContainer) Pod() vc.VCPod {
	return c.MockPod
}

func (c *VCMockContainer) Process() vc.Process {
	return c.MockProcess
}

func (c *VCMockContainer) GetToken() string {
	return c.MockToken
}

func (c *VCMockContainer) GetPid() int {
	return c.MockPid
}

func (c *VCMockContainer) SetPid(pid int) error {
	c.MockPid = pid
	return nil
}

func (c *VCMockContainer) URL() string {
	return c.MockURL
}

func (c *VCMockContainer) GetAnnotations() map[string]string {
	return c.MockAnnotations
}
