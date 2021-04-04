# docker.volume.driver.mount

Volume driver for Docker that can mount with mount cli.

## Install

```
└› docker plugin install sa4zet/docker.volume.driver.mount.plugin
Plugin "sa4zet/docker.volume.driver.mount.plugin" is requesting the following privileges:
 - network: [host]
 - capabilities: [CAP_SYS_ADMIN]
Do you grant the above permissions? [y/N] y
latest: Pulling from sa4zet/docker.volume.driver.mount.plugin
Digest: sha256:21b0ebeb5f242b5026b375c80741039273858b50309b27b544c9a86896d7e6a7
7a48fcf5abcd: Complete
Installed plugin sa4zet/docker.volume.driver.mount.plugin
```

## Check

```
└› docker plugin list
ID             NAME                                              DESCRIPTION           ENABLED
08d4c68a1917   sa4zet/docker.volume.driver.mount.plugin:latest   Mount volume driver   true
```

## Uninstall

```
└› docker plugin disable sa4zet/docker.volume.driver.mount.plugin
sa4zet/docker.volume.driver.mount.plugin

└› docker plugin rm sa4zet/docker.volume.driver.mount.plugin
sa4zet/docker.volume.driver.mount.plugin

└› docker plugin list
ID        NAME      DESCRIPTION   ENABLED
```

### Options

All available options are documented here and must be set.

|Key|Description|
|---|---|
|`fstype`|The filesystem types which are currently supported depend on the running kernel.|
|`source`|Mountable source|
|`options`|Filesystem type specified mount options as a comma-separated list.|

## Usage

### ad-hoc mount

```
└› docker run \
--rm \
-it \
--mount='type=volume,destination=/samba.mount/,"volume-driver=sa4zet/docker.volume.driver.mount.plugin","volume-opt=fstype=cifs","volume-opt=source=//samba/share","volume-opt=options=username=user.for.share,password=password.for.share,vers=3.0,dir_mode=0777,file_mode=0777,serverino"' \
debian:sid-slim \
/bin/bash

root@4cc06e15a839:/# ls /samba.mount/
shared.dir  shared.file
root@4cc06e15a839:/# exit
```

### Create volume

```
└› docker volume create \
--driver="sa4zet/docker.volume.driver.mount.plugin" \
--name="samba.mount" \
--opt="fstype=cifs" \
--opt="source=//samba/share" \
--opt="options=username=user.for.share,password=password.for.share,vers=3.0,dir_mode=0777,file_mode=0777,serverino" \
samba.mount

└› docker volume list
DRIVER                                            VOLUME NAME
sa4zet/docker.volume.driver.mount.plugin:latest   samba.mount

└› docker run \
--rm \
-it \
--mount="type=volume,source=samba.mount,destination=/samba.mount/" \
debian:sid-slim \
/bin/bash

root@59243434087d:/# ls /samba.mount/
shared.dir  shared.file

root@59243434087d:/# exit

└› docker volume rm samba.mount
samba.mount

└› docker volume list
DRIVER    VOLUME NAME
```

# License

https://github.com/sa4zet-org/docker.volume.driver.mount/blob/master/LICENSE
