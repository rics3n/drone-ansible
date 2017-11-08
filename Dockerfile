# Docker image for the Drone build runner
#
#     docker build --rm=true -t rics3n/drone-ansible .

FROM golang:1.9.2 AS builder
WORKDIR  /go/src/github.com/alexellis/drone-ansible/
COPY . .
RUN CGO_ENABLED=0 go build -a -tags netgo

FROM williamyeh/ansible:alpine3
COPY --from=builder /go/src/github.com/alexellis/drone-ansible/drone-ansible /bin/
ENTRYPOINT ["/bin/drone-ansible"]
