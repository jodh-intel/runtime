// Copyright (c) 2014,2015,2016 Docker, Inc.
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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	vc "github.com/containers/virtcontainers"
	"github.com/containers/virtcontainers/pkg/oci"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

type execParams struct {
	ociProcess   oci.CompatOCIProcess
	cID          string
	pidFile      string
	console      string
	consoleSock  string
	detach       bool
	processLabel string
	noSubreaper  bool
}

var execCLICommand = cli.Command{
	Name:  "exec",
	Usage: "Execute new process inside the container",
	ArgsUsage: `<container-id> <command> [command options]  || -p process.json <container-id>

   <container-id> is the name for the instance of the container and <command>
   is the command to be executed in the container. <command> can't be empty
   unless a "-p" flag provided.

EXAMPLE:
   If the container is configured to run the linux ps command the following
   will output a list of processes running in the container:

       # ` + name + ` <container-id> ps`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "console",
			Usage: "path to a pseudo terminal",
		},
		cli.StringFlag{
			Name:  "console-socket",
			Value: "",
			Usage: "path to an AF_UNIX socket which will receive a file descriptor referencing the master end of the console's pseudoterminal",
		},
		cli.StringFlag{
			Name:  "cwd",
			Usage: "current working directory in the container",
		},
		cli.StringSliceFlag{
			Name:  "env, e",
			Usage: "set environment variables",
		},
		cli.BoolFlag{
			Name:  "tty, t",
			Usage: "allocate a pseudo-TTY",
		},
		cli.StringFlag{
			Name:  "user, u",
			Usage: "UID (format: <uid>[:<gid>])",
		},
		cli.StringFlag{
			Name:  "process, p",
			Usage: "path to the process.json",
		},
		cli.BoolFlag{
			Name:  "detach,d",
			Usage: "detach from the container's process",
		},
		cli.StringFlag{
			Name:  "pid-file",
			Value: "",
			Usage: "specify the file to write the process id to",
		},
		cli.StringFlag{
			Name:  "process-label",
			Usage: "set the asm process label for the process commonly used with selinux",
		},
		cli.StringFlag{
			Name:  "apparmor",
			Usage: "set the apparmor profile for the process",
		},
		cli.BoolFlag{
			Name:  "no-new-privs",
			Usage: "set the no new privileges value for the process",
		},
		cli.StringSliceFlag{
			Name:  "cap, c",
			Value: &cli.StringSlice{},
			Usage: "add a capability to the bounding set for the process",
		},
		cli.BoolFlag{
			Name:   "no-subreaper",
			Usage:  "disable the use of the subreaper used to reap reparented processes",
			Hidden: true,
		},
	},
	Action: func(context *cli.Context) error {
		return execute(context)
	},
}

func generateExecParams(context *cli.Context, specProcess *oci.CompatOCIProcess) (execParams, error) {
	ctxArgs := context.Args()

	params := execParams{
		cID:          ctxArgs.First(),
		pidFile:      context.String("pid-file"),
		console:      context.String("console"),
		consoleSock:  context.String("console-socket"),
		detach:       context.Bool("detach"),
		processLabel: context.String("process-label"),
		noSubreaper:  context.Bool("no-subreaper"),
	}

	if context.IsSet("process") == true {
		var ociProcess oci.CompatOCIProcess

		fileContent, err := ioutil.ReadFile(context.String("process"))
		if err != nil {
			return execParams{}, err
		}

		if err := json.Unmarshal(fileContent, &ociProcess); err != nil {
			return execParams{}, err
		}

		params.ociProcess = ociProcess
	} else {
		params.ociProcess = *specProcess

		// Override terminal
		if context.IsSet("tty") {
			params.ociProcess.Terminal = context.Bool("tty")
		}

		// Override user
		if context.String("user") != "" {
			params.ociProcess.User = specs.User{
				Username: context.String("user"),
			}
		}

		// Override env
		params.ociProcess.Env = append(params.ociProcess.Env, context.StringSlice("env")...)

		// Override cwd
		if context.String("cwd") != "" {
			params.ociProcess.Cwd = context.String("cwd")
		}

		// Override no-new-privs
		if context.IsSet("no-new-privs") {
			params.ociProcess.NoNewPrivileges = context.Bool("no-new-privs")
		}

		// Override apparmor
		if context.String("apparmor") != "" {
			params.ociProcess.ApparmorProfile = context.String("apparmor")
		}

		params.ociProcess.Args = ctxArgs.Tail()
	}

	return params, nil
}

func execute(context *cli.Context) error {
	containerID := context.Args().First()
	status, podID, err := getExistingContainerInfo(containerID)
	if err != nil {
		return err
	}

	// Retrieve OCI spec configuration.
	ociSpec, err := oci.GetOCIConfig(status)
	if err != nil {
		return err
	}

	params, err := generateExecParams(context, ociSpec.Process)
	if err != nil {
		return err
	}

	params.cID = status.ID

	// container MUST be running
	if status.State.State != vc.StateRunning {
		return fmt.Errorf("Container %s is not running", params.cID)
	}

	envVars, err := oci.EnvVars(params.ociProcess.Env)
	if err != nil {
		return err
	}

	consolePath, err := setupConsole(params.console, params.consoleSock)
	if err != nil {
		return err
	}

	cmd := vc.Cmd{
		Args:        params.ociProcess.Args,
		Envs:        envVars,
		WorkDir:     params.ociProcess.Cwd,
		User:        params.ociProcess.User.Username,
		Interactive: params.ociProcess.Terminal,
		Console:     consolePath,
		Detach:      noNeedForOutput(params.detach, params.ociProcess.Terminal),
	}

	_, _, process, err := vc.EnterContainer(podID, params.cID, cmd)
	if err != nil {
		return err
	}

	// Creation of PID file has to be the last thing done in the exec
	// because containerd considers the exec to have finished starting
	// after this file is created.
	if err := createPIDFile(params.pidFile, process.Pid); err != nil {
		return err
	}

	if !params.detach {
		p, err := os.FindProcess(process.Pid)
		if err != nil {
			return err
		}

		ps, err := p.Wait()
		if err != nil {
			return fmt.Errorf("Process state %s, container info %+v: %v",
				ps.String(), status, err)
		}

		// Exit code has to be forwarded in this case.
		return cli.NewExitError("", ps.Sys().(syscall.WaitStatus).ExitStatus())
	}

	return nil
}
