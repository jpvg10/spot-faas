package main

import (
	"bytes"
	"log"
	"os/exec"
)

func createVM(name string) string {
	return run("gcloud",
		"compute",
		"instances",
		"create",
		name,
		"--provisioning-model=SPOT",
		"--instance-termination-action=DELETE",
		"--format=json",
	)
}

func deleteVM(name string) {
	run("gcloud",
		"compute",
		"instances",
		"delete",
		name,
	)
}

func run(command string, args ...string) string {
	cmd := exec.Command(command, args...)

	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	err := cmd.Run()

	if err != nil {
		log.Print(cmdErr.String())
		log.Fatal(err)
	}

	log.Printf("Command output: %s\n", cmdOut.String())
	return cmdOut.String()
}
