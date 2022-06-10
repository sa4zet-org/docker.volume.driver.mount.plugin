[![Check and build](https://github.com/sa4zet-org/docker.volume.driver.mount.plugin/actions/workflows/release.yml/badge.svg)](https://github.com/sa4zet-org/docker.volume.driver.mount.plugin/actions/workflows/release.yml)

# docker.volume.driver.mount.plugin

Volume driver for Docker that can mount with mount cli.

## Install

```bash
docker plugin install --grant-all-permissions ghcr.io/sa4zet-org/docker.volume.driver.mount.plugin
```

## Uninstall

```bash
docker plugin disable ghcr.io/sa4zet-org/docker.volume.driver.mount.plugin
```
```bash
docker plugin rm ghcr.io/sa4zet-org/docker.volume.driver.mount.plugin
```

### Options

All available options are documented here and must be set.

| Key       | Description                                                                      |
|-----------|----------------------------------------------------------------------------------|
| `fstype`  | The filesystem types which are currently supported depend on the running kernel. |
| `source`  | Mountable source                                                                 |
| `options` | Filesystem type specified mount options as a comma-separated list.               |

## Usage

### ad-hoc mount

```bash
docker run \
--rm \
-it \
--mount='type=volume,destination=/samba.mount/,"volume-driver=ghcr.io/sa4zet-org/docker.volume.driver.mount.plugin","volume-opt=fstype=cifs","volume-opt=source=//samba/share","volume-opt=options=username=user.for.share,password=password.for.share,vers=3.0,dir_mode=0777,file_mode=0777,serverino"' \
debian:sid-slim \
/bin/bash
```

### Create volume

```bash
docker volume create \
--driver="ghcr.io/sa4zet-org/docker.volume.driver.mount.plugin" \
--name="samba.mount" \
--opt="fstype=cifs" \
--opt="source=//samba/share" \
--opt="options=username=user.for.share,password=password.for.share,vers=3.0,dir_mode=0777,file_mode=0777,serverino" \
samba.mount
```
```bash
docker run \
--rm \
-it \
--mount="type=volume,source=samba.mount,destination=/samba.mount/" \
debian:sid-slim \
/bin/bash
```

# License

https://github.com/sa4zet-org/docker.volume.driver.mount.plugin/blob/master/LICENSE
