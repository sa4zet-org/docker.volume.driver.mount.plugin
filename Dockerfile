FROM golang:1.16.2-alpine3.12 as builder

WORKDIR github.com/sa4zet-org/docker.volume.driver.mount
COPY src/ .

RUN apk update \
  && apk add git \
  && go build \
    --ldflags "-extldflags -static" \
    -o /usr/bin/mountVolumeDriver

FROM alpine:3.12
COPY --from=builder /usr/bin/mountVolumeDriver /usr/bin/mountVolumeDriver
