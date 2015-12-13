# drone-docker
Drone plugin for provisioning with ansible

## Usage

```
./drone-ansible <<EOF
{
    "build": {
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src"
    },
    "vargs": {
        "playbook": "provisioning/provision.yml",
		"inventory": "provisioning/inventory/staging"
    }
}
EOF

## Docker

Build the Docker container:

```sh
docker build --rm=true -t rics3n/drone-ansible .
```

Build and Publish a Docker container

```sh
docker run -i --privileged -v $(pwd):/drone/src rics3n/drone-ansible <<EOF
{
	"workspace": {
	 	"root": "/drone/src",
		"path": "/drone/src"
	},
	"build": {
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
	"vargs": {
		"playbook": "provisioning/provision.yml",
		"inventory": "provisioning/inventory/staging"
	}
}
EOF
```
