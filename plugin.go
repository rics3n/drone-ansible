package main

import (
	"fmt"
	"io"
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
		InventoryPath     string
		Inventories       []string
		Playbook          string
		SSHKey            string
		SSHKeyFile        string
		AnsibleConfigFile string
	}

	// Plugin defines the Ansible plugin parameters.
	Plugin struct {
		Build  Build
		Config Config // Ansible config
	}
)

func (p Plugin) Exec() error {
	// write the rsa private key
	privateKeyFile, err := writeKey(p.Config)
	if err != nil {
		return err
	}

	// write ansible configuration
	if err := writeAnsibleConf(p.Config); err != nil {
		return err
	}
	var cmds []*exec.Cmd
	cmds = append(cmds, commandVersion())

	for _, inventory := range p.Config.Inventories {
		cmds = append(cmds, command(p.Build, p.Config, inventory, privateKeyFile)) // docker tag
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

func command(build Build, config Config, inventory string, privateKeyFile string) *exec.Cmd {

	args := []string{
		commandEnvVars(build, privateKeyFile),
		"-i",
		filepath.Join(build.Path, config.InventoryPath, inventory),
		filepath.Join(build.Path, config.Playbook),
	}
	return exec.Command(ansibleBin, args...)
}

// helper function to create the docker info command.
func commandVersion() *exec.Cmd {
	return exec.Command(ansibleBin, "--version")
}

func commandEnvVars(build Build, privateKeyFile string) string {
	args := []string{
		fmt.Sprintf("-e ansible_ssh_private_key_file=%s", privateKeyFile),
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
func writeKey(config Config) (string, error) {
	home := "/root"
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	}
	sshpath := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshpath, 0700); err != nil {
		return "", err
	}
	confpath := filepath.Join(sshpath, "config")

	ioutil.WriteFile(confpath, []byte("StrictHostKeyChecking no\n"), 0700)

	privpath := filepath.Join(sshpath, "id_rsa")
	if len(config.SSHKey) != 0 {
		ioutil.WriteFile(privpath, []byte(config.SSHKey), 0600)
		return privpath, nil
	}

	if len(config.SSHKeyFile) != 0 {
		err := copyFile(config.SSHKeyFile, privpath)
		if err != nil {
			serr := err.Error()
			return "", fmt.Errorf("Couldn't copy key file %s to %s with error: %q\n", config.SSHKeyFile, privpath, serr)
		}

		// Set correct permissions
		err = os.Chmod(privpath, 0400)
		if err != nil {
			serr := err.Error()
			return privpath, fmt.Errorf("Couldn't set permissions for file %s with error: %q\n", privpath, serr)
		}

		return privpath, nil
	}

	return "", nil
}

func writeAnsibleConf(config Config) error {
	confpath := "/etc/ansible/ansible.cfg"

	if len(config.AnsibleConfigFile) != 0 {
		err := copyFile(config.AnsibleConfigFile, confpath)
		if err != nil {
			return err
		}
	} else {
		f, err := os.OpenFile(confpath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		//this disables host key checking.. be aware of the man in the middle
		if _, err := f.Write([]byte("[defaults]\nhost_key_checking = False\n")); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Copies a file
func copyFile(src, dst string) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Couldn't stat source file %s", src)
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}

	dfi, err := os.Stat(dst)
	if os.SameFile(sfi, dfi) {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Couldn't open source file %s", src)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("Couldn't create destination file %s", src)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("Couldn't copy file from source %s to destination %s", src, dst)
	}
	err = out.Sync()
	return nil

}
