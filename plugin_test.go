package main

import (
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	// without user
	if command := strings.Join(command(Build{
		Path: "path/to/playbookdir",
		SHA:  "123456789ABCDEF",
		Tag:  "TAGNAME",
	}, Config{
		InventoryPath: "inventory-path",
		Inventories:   []string{"inventory1", "inventory2"},
		Playbook:      "playbook.yml",
		SSHKey:        "",
	}, "inventory").Args, " "); command != "/usr/bin/ansible-playbook -e ansible_ssh_private_key_file=/root/.ssh/id_rsa -e commit_sha=123456789ABCDEF -e commit_tag=TAGNAME -i path/to/playbookdir/inventory-path/inventory path/to/playbookdir/playbook.yml" {
		t.Error("unexpected command: ", command)
	}

	// with user
	if command := strings.Join(command(Build{
		Path: "path/to/playbookdir",
		SHA:  "123456789ABCDEF",
		Tag:  "TAGNAME",
	}, Config{
		InventoryPath: "inventory-path",
		Inventories:   []string{"inventory1", "inventory2"},
		Playbook:      "playbook.yml",
		RemoteUser:    "user1",
		SSHKey:        "",
	}, "inventory").Args, " "); command != "/usr/bin/ansible-playbook -e ansible_ssh_private_key_file=/root/.ssh/id_rsa -e commit_sha=123456789ABCDEF -e commit_tag=TAGNAME -u user1 -i path/to/playbookdir/inventory-path/inventory path/to/playbookdir/playbook.yml" {
		t.Error("unexpected command: ", command)
	}
}
