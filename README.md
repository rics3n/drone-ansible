# drone-docker
Drone plugin for provisioning with ansible


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
		"path": "/drone/src"
	},
	"build" : {
		"number": 1,
		"head_commit": {
			"sha": "9f2849d5",
			"branch": "master",
			"ref": "refs/heads/master"
		}
	},
	"vargs": {
		"playbook": "provision.yml",
		"inventory": "inventory/staging"
	}
}
EOF
```
