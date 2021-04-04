package main

import (
	"github.com/docker/go-plugins-helpers/volume"
)

func main() {
	volumeHandler := volume.NewHandler(newMountVolumeDriver())
	if err := volumeHandler.ServeUnix("mountVolumeDriver", 0); err != nil {
		panic(err)
	}
}
