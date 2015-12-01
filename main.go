package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/drone/drone-plugin-go/plugin"
)

type Ansible struct {
	Inventory 	string   `json:"inventory"`
	Playbook 	string   `json:"playbook"`
}

func main() {
	workspace := plugin.Workspace{}
	build := plugin.Build{}
	vargs := Ansible{}

	plugin.Param("workspace", &workspace)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	// Set the Inventory file
	if len(vargs.Inventory) == 0 {
		vargs.Inventory = "hosts"
	}
	// Set the Playbook
	if len(vargs.Playbook) == 0 {
		vargs.Playbook = "playbook.yml"
	}

	// write the rsa private key
	if err := writeKey(workspace); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Docker environment info
	cmd := exec.Command("/usr/bin/ansible-playbook", "-i", vargs.Inventory, vargs.Playbook)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}

// Writes the RSA private key
func writeKey(in plugin.Workspace) error {
	if in.Keys == nil || len(in.Keys.Private) == 0 {
		return nil
	}
	home := "/root"
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	}
	sshpath := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshpath, 0700); err != nil {
		return err
	}
	confpath := filepath.Join(sshpath, "config")
	privpath := filepath.Join(sshpath, "id_rsa")
	ioutil.WriteFile(confpath, []byte("StrictHostKeyChecking no\n"), 0700)
	return ioutil.WriteFile(privpath, []byte(in.Keys.Private), 0600)
}