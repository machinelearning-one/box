package utils

import (
	"fmt"

	"github.com/docker/docker/api/types"
)

func TruePath(cmp string, vols []types.MountPoint) (string, error) {
	match := -1
	dest := ""
	for i, pair := range vols {
		dest = pair.Destination
		if dest == cmp[:len(dest)] {
			match = i
			break
		}
	}
	if match == -1 {
		return "", fmt.Errorf("could not find matching mount point for %s", cmp)
	} else {
		return vols[match].Source + cmp[len(dest):], nil
	}
}
