Use the Docker plugin to provision with ansible.
The following parameters are used to configure this plugin:

* `inventory` - define the inventory file (default: hosts)
* `playbook` - define the playbook file (default: playbook.yml)

The following is a sample Docker configuration in your .drone.yml file:

```yaml
provision:
  ansible:
    inventory: inventory/staging
    playbook: provision.yml
```
