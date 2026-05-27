package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Ensures gofmt doesn't remove the imports above (feel free to remove this!)
var _ = os.Args
var _ = exec.Command

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {

	// fmt.Println("os.Args[0]: ", os.Args[0])
	// fmt.Println("os.Args[1]: ", os.Args[1])
	// fmt.Println("os.Args[2]: ", os.Args[2])
	// fmt.Println("os.Args[3:]: ", os.Args[3:])
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(output))
}
