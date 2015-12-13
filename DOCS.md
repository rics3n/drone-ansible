Use the Drone plugin to provision with ansible.
The following parameters are used to configure this plugin:

* `inventory` - define the inventory file (default: hosts)
* `playbook` - define the playbook file (default: playbook.yml)

The following is a sample Docker configuration in your .drone.yml file:

```yaml
deploy:
  ansible:
  	image: rics3n/drone-ansible
    inventory: inventory/staging
    playbook: provision.yml
```
