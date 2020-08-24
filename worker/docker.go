package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ohsu-comp-bio/funnel/events"
)

// DockerCommand is responsible for configuring and running a docker container.
type DockerCommand struct {
	Image           string
	Command         []string
	Volumes         []Volume
	Workdir         string
	ContainerName   string
	RemoveContainer bool
	Env             map[string]string
	Stdin           io.Reader
	Stdout          io.Writer
	Stderr          io.Writer
	Event           *events.ExecutorWriter
}

// Run runs the Docker command and blocks until done.
func (dcmd DockerCommand) Run(ctx context.Context) error {
	// Sync docker API version info.
	err := SyncDockerAPIVersion()
	if err != nil {
		dcmd.Event.Error("failed to sync docker client API version", err)
	}

	pullcmd := exec.Command("ls", dcmd.Image)
	err = pullcmd.Run()
	if err != nil {
		dcmd.Event.Error("failed to locate singularity image", err)
	}

	args := []string{"exec"}

	// if dcmd.Env != nil {
	// 	for k, v := range dcmd.Env {
	// 		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	// 	}
	// }

	if dcmd.ContainerName != "" {
		dcmd.Event.Error("Container name not supported")
	}

	if dcmd.Workdir != "" {
		args = append(args, "--pwd", dcmd.Workdir)
	}

	for _, vol := range dcmd.Volumes {
		arg := formatVolumeArg(vol, dcmd)
		args = append(args, "--bind", arg)
	}

	args = append(args, dcmd.Image)
	args = append(args, dcmd.Command...)

	// Roughly: `docker run --rm -i --read-only -w [workdir] -v [bindings] [imageName] [cmd]`
	dcmd.Event.Info("Running command", "cmd", "singularity "+strings.Join(args, " "))
	cmd := exec.Command("singularity", args...)

	if dcmd.Stdin != nil {
		cmd.Stdin = dcmd.Stdin
	}
	if dcmd.Stdout != nil {
		cmd.Stdout = dcmd.Stdout
	}
	if dcmd.Stderr != nil {
		cmd.Stderr = dcmd.Stderr
	}
	go dcmd.inspectContainer(ctx)
	return cmd.Run()
}

// Stop stops the container.
func (dcmd DockerCommand) Stop() error {
	dcmd.Event.Info("Stopping container", "container", dcmd.ContainerName)
	cmd := exec.Command("docker", "stop", dcmd.ContainerName)
	return cmd.Run()
}

func formatVolumeArg(v Volume, dcmd DockerCommand) string {
	// `o` is structed as "HostPath:ContainerPath".
	if v.Readonly {
		dcmd.Event.Error("Read only mount not supported")
	}
	return fmt.Sprintf("%s:%s", v.HostPath, v.ContainerPath)
}

type metadata struct {
	ID    string
	Name  string
	Image string
}

// inspectContainer inspects the docker container for metadata.
func (dcmd *DockerCommand) inspectContainer(ctx context.Context) {
	// Give the container time to start.
	time.Sleep(2 * time.Second)

	// Inspect the container for metadata
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cmd := exec.CommandContext(ctx, "docker", "inspect", dcmd.ContainerName)
			out, err := cmd.Output()
			if err == nil {
				meta := []metadata{}
				err := json.Unmarshal(out, &meta)
				if err == nil && len(meta) == 1 {
					dcmd.Event.Info("container metadata",
						"containerID", meta[0].ID,
						"containerName", meta[0].Name,
						"containerImageHash", meta[0].Image)
					return
				}
			}
		}
	}
}

type dockerVersion struct {
	Client string
	Server string
}

// SyncDockerAPIVersion ensures that the client uses the same API version as
// the server.
func SyncDockerAPIVersion() error {
	if os.Getenv("DOCKER_API_VERSION") == "" {
		cmd := exec.Command("docker", "version", "--format", `{"Server": "{{.Server.APIVersion}}", "Client": "{{.Client.APIVersion}}"}`)
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("docker version command failed: %v", err)
		}
		version := &dockerVersion{}
		err = json.Unmarshal(out, version)
		if err != nil {
			return fmt.Errorf("failed to unmarshal docker version: %v", err)
		}
		if version.Client != version.Server {
			os.Setenv("DOCKER_API_VERSION", version.Server)
		}
	}
	return nil
}
