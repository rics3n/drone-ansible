package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	defaultPlaybook      = "provisioning/provision.yml"
	defaultInventoryPath = "provisioning/inventory"
	ansibleBin           = "/usr/bin/ansible-playbook"
)

type (
	//Config defined the ansible configuration params
	Build struct {
		Path string
		SHA  string
		Tag  string
	}

	//Config defined the ansible configuration params
	Config struct {
		InventoryPath string
		Inventories   []string
		Playbook      string
		SSHKey        string
		SSHUser       string
	}

	// Plugin defines the Ansible plugin parameters.
	Plugin struct {
		Build  Build
		Config Config // Ansible config
	}
)

func (p Plugin) Exec() error {
	// write the rsa private key
	if err := writeKey(p.Config); err != nil {
		return err
	}

	// write ansible configuration
	if err := writeAnsibleConf(); err != nil {
		return err
	}
	var cmds []*exec.Cmd
	cmds = append(cmds, commandVersion())

	for _, inventory := range p.Config.Inventories {
		cmds = append(cmds, command(p.Build, p.Config, inventory)) // docker tag
	}

	// Run ansible
	// execute all commands in batch mode.
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func command(build Build, config Config, inventory string) *exec.Cmd {
	var args []string
	args = append(args, commandEnvVars(build))
	if config.SSHUser != "" {
		args = append(args, "-u", config.SSHUser)
	}
	args = append(args, "-i", filepath.Join(build.Path, config.InventoryPath, inventory))
	args = append(args, filepath.Join(build.Path, config.Playbook))
	return exec.Command(ansibleBin, args...)
}

// helper function to create the docker info command.
func commandVersion() *exec.Cmd {
	return exec.Command(ansibleBin, "--version")
}

func commandEnvVars(build Build) string {
	args := []string{
		"-e ansible_ssh_private_key_file=/root/.ssh/id_rsa",
		fmt.Sprintf("-e commit_sha=%s", build.SHA),
	}

	if len(build.Tag) != 0 {
		args = append(args, fmt.Sprintf("-e commit_tag=%s", build.Tag))
	}

	return strings.Join(args, " ")
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}

// Writes the RSA private key
func writeKey(config Config) error {
	if len(config.SSHKey) == 0 {
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
	return ioutil.WriteFile(privpath, []byte(config.SSHKey), 0600)
}

func writeAnsibleConf() error {
	confpath := "/etc/ansible/ansible.cfg"
	//this disables host key checking.. be aware of the man in the middle
	return ioutil.WriteFile(confpath, []byte("[defaults]\nhost_key_checking = False\n"), 0600)
}
