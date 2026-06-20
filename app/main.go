//go:build linux

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Usage: your_docker.sh run <image> <command path> <arg1> <arg2> ...
func main() {

	commandPath := os.Args[3]
	args := os.Args[4:]

	rootDir, err := os.MkdirTemp("", "mydocker-*")
	// fmt.Printf("Created root dir: %s\n", rootDir)
	if err != nil {
		fmt.Printf("Err creating temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(rootDir)

	// file system isolation
	if err := copyBinary(commandPath, rootDir); err != nil {
		fmt.Printf("Err copying binary: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Chroot(rootDir); err != nil {
		fmt.Printf("Err chroot: %v\n", err)
		os.Exit(1)
	}
	if err := os.Chdir("/"); err != nil {
		fmt.Printf("Err chdir: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(commandPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// process isolation
	// Cloneflags and CLONE_NEWPID are Linux-only fields/constants on syscall.SysProcAttr.
	// Editing/building on macOS uses the darwin syscall package, which defines neither,
	// so gopls/go build will flag them as undefined there. This file is meant to be
	// compiled inside the Linux container (per your_docker.sh / Dockerfile), where
	// GOOS=linux and both symbols exist.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Err: %v", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}

}

func copyBinary(src string, rootDir string) error {
	destPath := filepath.Join(rootDir, src)
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return os.Chmod(destPath, 0755)
}
