package utils

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func InsideContainer() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

func ContainerIdentity() (string, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/self/mountinfo | grep '/docker/containers/' | head -1 | awk '{print $4}' | cut -d '/' -f6")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String()[:12], nil
}

func VolumeMappings(ctx context.Context, client *client.Client, containerID string) ([]types.MountPoint, error) {
	container, err := client.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	return container.Mounts, nil
}
