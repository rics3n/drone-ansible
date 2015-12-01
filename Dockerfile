# Docker image for the ansible plugin
#
#     docker build --rm=true -t rics3n/drone-ansible .

FROM williamyeh/ansible:alpine3

RUN apk add -U ca-certificates openssh && rm -rf /var/cache/apk/*

ADD drone-ansible /bin/
ENTRYPOINT ["/bin/drone-ansible"]