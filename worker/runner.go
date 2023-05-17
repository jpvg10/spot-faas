package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func runJob(param string) string {
	dockerCommand := []string{"run"}

	if len(param) > 0 {
		dockerCommand = append(dockerCommand, "-e", fmt.Sprintf("MESSAGE=%v", param))
	}
	dockerCommand = append(dockerCommand, "worker")

	cmd := exec.Command("docker", dockerCommand...)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	err := cmd.Run()

	if err != nil {
		log.Print(cmdErr.String())
		log.Fatal(err)
	}

	log.Printf("Container output: %s\n", cmdOut.String())
	return cmdOut.String()
}
