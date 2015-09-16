package hypercore

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/docker/machine/libmachine/log"
)

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reEqualLine       = regexp.MustCompile(`(.+)=(.*)`)
	reEqualQuoteLine  = regexp.MustCompile(`"(.+)"="(.*)"`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

var (
	ErrMachineExist      = errors.New("machine already exists")
	ErrMachineNotExist   = errors.New("machine does not exist")
	ErrHypercoreNotFound = errors.New("hypercore binary (linux) not found")
	HypercoreCmd         = setHypercoreCmd()
)

// Detect the hypercore(cmdname: linux) cmd's path if needed
func setHypercoreCmd() string {
	cmd := "linux"
	if path, err := exec.LookPath(cmd); err == nil {
		return path
	}
	return cmd
}

func hypercoreExec(args ...string) (string, error) {
	var Password string
	cmd := exec.Command("sudo", HypercoreCmd, strings.Join(args, " "))
	cmd.Stdin = strings.NewReader(Password)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	log.Debugf("executing: %v %v %v", cmd, args, stdout.String())

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(stdout.String())
	return stdout.String(), err
}

func hypercoreInit() {
	hypercoreExec("init")
}

func hypercoreBoot() {
	hypercoreExec("boot")
}

func hypercoreSsh() {
	hypercoreExec("ssh")
}

func hypercoroOut(args ...string) (string, error) { // TODO
	stdout, _, err := hypercoreOutErr(args...)
	return stdout, err
}

func hypercoreOutErr(args ...string) (string, string, error) { // TODO
	cmd := exec.Command(HypercoreCmd, args...)
	log.Debugf("executing: %v %v", HypercoreCmd, strings.Join(args, " "))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	stderrStr := stderr.String()
	log.Debugf("STDOUT: %v", stdout.String())
	log.Debugf("STDERR: %v", stderrStr)
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrHypercoreNotFound
		}
	} else {
		if strings.Contains(stderrStr, "error:") {
			err = fmt.Errorf("%v %v failed: %v", HypercoreCmd, strings.Join(args, " "), stderrStr)
		}
	}
	return stdout.String(), stderrStr, err
}
