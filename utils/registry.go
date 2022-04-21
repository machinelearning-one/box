package utils

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func AvailableImages(ctx context.Context, client *client.Client) ([]string, error) {
	images, err := client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	var imageNames []string
	for _, image := range images {
		if image.RepoTags[0] != "<none>:<none>" {
			imageNames = append(imageNames, image.RepoTags...)
		}
	}
	return imageNames, nil
}

func RemoveImage(image string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	_, err = cli.ImageRemove(ctx, image, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}
