#!/bin/bash

docker plugin disable sa4zet/docker.volume.driver.mount.plugin
docker plugin rm sa4zet/docker.volume.driver.mount.plugin
docker build . --tag sa4zet/docker.volume.driver.mount
rm -rf rootfs
mkdir -p rootfs/
id=$(docker create sa4zet/docker.volume.driver.mount true)
docker export "${id}" | tar -x -C rootfs
docker rm -f "${id}"
docker plugin create sa4zet/docker.volume.driver.mount.plugin .
docker plugin enable sa4zet/docker.volume.driver.mount.plugin
