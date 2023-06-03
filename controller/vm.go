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
		"--zone=europe-north1-c",
		"--provisioning-model=SPOT",
		"--instance-termination-action=DELETE",
		"--image-project=spot-380110",
		"--image=worker-base-image",
		"--scopes=storage-ro",
		`--metadata=startup-script=#! /bin/bash
		mkdir /program
		gcloud storage cp gs://spot-thesis-files-2994/worker /program/worker
		chmod +x /program/worker
		/program/worker >> /program/log 2>&1`,
		"--format=json",
	)

	var data CreateResponse
	err := json.Unmarshal([]byte(output), &data)
	if err != nil {
		log.Printf("Could not unmarshal json: %s\n", err)
		return ""
	}

	return data[0].NetworkInterfaces[0].NetworkIP
}

func deleteVM(name string) {
	run("gcloud",
		"compute",
		"instances",
		"delete",
		name,
		"--zone=europe-north1-c",
		"--quiet",
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
