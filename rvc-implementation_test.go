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

// Description: Switch to the mock implementation of virtcontainers for
// testing..

package main

import (
	"fmt"

	"github.com/containers/virtcontainers/pkg/vcMock"
)

// testingImpl is a concrete mock RVC implementation used for testing
var testingImpl = &vcMock.VCMock{}

func init() {
	fmt.Printf("INFO: switching to fake virtcontainers implementation for testing\n")
	vci = testingImpl
}
