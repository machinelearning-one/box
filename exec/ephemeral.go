package exec

import (
	"box/core"
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func Ephemeral(ctx context.Context, client *client.Client, box core.Box) error {
	containerConfig := &container.Config{
		Image: box.Image,
		Cmd:   box.Cmd,
	}
	if box.WorkDir != "" {
		containerConfig.WorkingDir = box.WorkDir
	}
	hostConfig := &container.HostConfig{}
	if box.Vol.Target != "" {
		hostConfig.Mounts = []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: box.Vol.Source,
				Target: box.Vol.Target,
			},
		}
	}
	if box.GPUs {
		hostConfig.DeviceRequests = []container.DeviceRequest{
			{
				Count:  -1,
				Driver: "nvidia",
				// https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/user-guide.html#driver-capabilities
				Capabilities: [][]string{
					{
						"compute",
						"utility",
					},
				},
			},
		}
	}
	resp, err := client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return err
	}
	if err := client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	statusCh, errCh := client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return err
	case <-statusCh:
	}
	out, err := client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	return client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
}
