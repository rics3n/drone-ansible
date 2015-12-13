# drone-docker
Drone plugin for provisioning with ansible

## Docker

Build and Publish a Docker container

```sh
docker build --rm=true -t rics3n/drone-ansible .
docker push rics3n/drone-ansible
```

Run the docker container 

```sh
docker run -i --privileged -v $(pwd):/drone/src rics3n/drone-ansible <<EOF
{
	"workspace": {
	 	"root": "/drone/src",
		"path": "/drone/src",
		"keys": {
			"private": "..."
		}
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
