package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func runJob(params string) string {
	cmd := exec.Command("docker", "run", "worker")

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	err := cmd.Run()

	if err != nil {
		fmt.Println(cmdErr.String())
		log.Fatal(err)
	}

	fmt.Printf("Container output: %s\n", cmdOut.String())
	return cmdOut.String()
}
