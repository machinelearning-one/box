package utils

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func TouchBoxFile() error {
	// Check if .box file exists in current directory
	_, err := os.Stat(".box")
	if err == nil {
		return nil
	}
	// If not, create it
	f, err := os.Create(".box")
	if err != nil {
		return err
	}
	defer f.Close()
	// Write empty json to file
	_, err = f.Write([]byte("{}"))
	if err != nil {
		return err
	}
	return nil
}

func TouchImage(check string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	availableImages, err := AvailableImages(ctx, cli)
	if err != nil {
		return err
	}
	// Check if image is already present
	for _, image := range availableImages {
		if image == check {
			return nil
		}
	}
	// If not, pull it
	out, err := cli.ImagePull(ctx, check, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)
	return nil
}
