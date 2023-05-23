package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

func createVM(name string) string {
	output := run("gcloud",
		"compute",
		"instances",
		"create",
		name,
		"--provisioning-model=SPOT",
		"--instance-termination-action=DELETE",
		"--image-family=ubuntu-minimal-2204-lts",
		"--image-project=ubuntu-os-cloud",
		"--format=json",
	)

	var data CreateResponse
	err := json.Unmarshal([]byte(output), &data)
	if err != nil {
		log.Printf("Could not unmarshal json: %s\n", err)
		return ""
	}

	log.Println(data[0].NetworkInterfaces[0].NetworkIP)

	return data[0].NetworkInterfaces[0].NetworkIP
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

	// log.Printf("Command output: %s\n", cmdOut.String())
	return cmdOut.String()
}
