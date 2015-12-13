# Docker image for the Drone build runner
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t rics3n/drone-ansible .

#FROM ubuntu:14.04

#RUN apt-get update && \
#    apt-get install --no-install-recommends -y software-properties-common ca-certificates && \
#    apt-add-repository ppa:ansible/ansible && \
#    apt-get update && \
#    apt-get install -y ansible && \
#    rm -rf /var/lib/apt/lists/*

FROM williamyeh/ansible:alpine3

ADD drone-ansible /bin/

ENTRYPOINT ["/bin/drone-ansible"]