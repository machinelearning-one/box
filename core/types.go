package core

import (
	"box/utils"
	"context"
	"os"

	"github.com/docker/docker/client"
)

type Vol struct {
	Source string
	Target string
}

type Box struct {
	Image   string
	Cmd     []string
	Vol     Vol
	WorkDir string
	GPUs    bool
}

type BoxConfig struct {
	Image   string   `json:"image"`
	Cmd     []string `json:"cmd"`
	Target  string   `json:"target"`
	WorkDir string   `json:"workdir"`
	GPUs    bool     `json:"gpus"`
}

func ConfigToBox(config BoxConfig, rest []string) (Box, error) {
	// Add the rest of the arguments to the command
	cmd := append(config.Cmd, rest...)
	box := Box{
		Image:   config.Image,
		Cmd:     cmd,
		WorkDir: config.WorkDir,
		GPUs:    config.GPUs,
	}
	if config.Target != "" {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return box, err
		}
		path := cwd
		// Check if inside a container
		if utils.InsideContainer() {
			// Get container ID
			containerID, err := utils.ContainerIdentity()
			if err != nil {
				return box, err
			}
			// Get volume mappings
			client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				return box, err
			}
			mappings, err := utils.VolumeMappings(context.Background(), client, containerID)
			if err != nil {
				return box, err
			}
			// Get the true path
			path, err = utils.TruePath(cwd, mappings)
			if err != nil {
				return box, err
			}
		}

		box.Vol = Vol{
			Source: path,
			Target: config.Target,
		}
	}
	return box, nil
}
