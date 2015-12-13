# Docker image for the Drone build runner
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t rics3n/drone-ansible .

FROM alpine:3.2

RUN apk update && \
  	apk add \
    	ca-certificates \
    	ansible \
  rm -rf /var/cache/apk/*

ADD drone-ansible /bin/

ENTRYPOINT ["/bin/drone-ansible"]