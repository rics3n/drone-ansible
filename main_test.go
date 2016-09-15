package main

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/drone/drone-plugin-go/plugin"
)

func TestCommandBuildCorrectly(t *testing.T) {
	vargs := Ansible{}
	vargs.Playbook = "test.yml"
	vargs.Inventory = "inventory/staging"
	w := plugin.Workspace{Path: "/test/path"}
	if !reflect.DeepEqual(command(vargs, w).Args, []string{
		"/usr/bin/ansible-playbook",
		"-i",
		filepath.Join(w.Path, vargs.Inventory),
		filepath.Join(w.Path, vargs.Playbook),
	}) {
		trace(command(vargs, w))
		t.Error("command not composed correctly")
	}
}
