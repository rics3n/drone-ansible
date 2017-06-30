# drone-ansible

Use the Drone plugin to provision with ansible.
The following parameters are used to configure this plugin:

* `inventory` - define the inventory file (default: provisioning/inventory/staging)
* `inventories` - define multiple inventory files to deploy
* `playbook` - define the playbook file (default: provisioning/provision.yml)
* `ssh-key` - define the ssh-key to use for connecting to hosts

The following is a sample Docker configuration in your .drone.yml file:

```yaml
pipeline:
  deploy-staging:
  	image: rics3n/drone-ansible
    inventory: inventory/staging
    playbook: provision.yml
	secrets: [ ssh_key ]
	when:
		branch: master
```

```yaml
pipeline:
  deploy-staging:
  	image: rics3n/drone-ansible
    inventories: [ staging, latest ]
    playbook: provision.yml
	secrets: [ ssh_key ]
	when:
		branch: master
```

To addthe ssh key use drone secrets via the cli

```
drone secret add \
  -repository user/repo \
  -image rics3n/drone-ansible \
  -name ssh_key \
  -value @Path/to/.ssh/id_rsa
```

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t rics3n/drone-ansible:2 .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-ansible' not found or does not exist..
```

## Usage

Execute from a project directory:

```
docker run --rm=true \
  -e PLUGIN_SSH_KEY=${SSH_KEY} \
  -e DRONE_WORKSPACE=/go/src/github.com/username/test \
  -e PLUGIN_INVENTORY=provisioning/inventory/homerlatest \
  -v $(pwd):/go/src/github.com/username/test \
  -w /go/src/github.com/rics3n/test \
  rics3n/drone-ansible:2
```
