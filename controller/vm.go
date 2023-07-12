package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

func createVM(name string) (string, error) {
	cmdOut, cmdErr := run("gcloud",
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
		"--machine-type=e2-standard-2",
		"--metadata",
		`startup-script=#! /bin/bash
		mkdir /program
		gcloud storage cp gs://spot-thesis-files-2994/worker /program/worker
		chmod +x /program/worker
		/program/worker >> /program/log 2>&1

		,shutdown-script=#! /bin/bash
		pid=$(pidof worker)
		kill -s SIGTERM $pid`,
		"--format=json",
	)

	if cmdErr != nil {
		return "", cmdErr
	}

	var data CreateResponse
	err := json.Unmarshal([]byte(cmdOut), &data)
	if err != nil {
		log.Printf("Could not unmarshal json: %s", err)
		return "", err
	}

	return data[0].NetworkInterfaces[0].NetworkIP, nil
}

func deleteVM(name string) error {
	_, err := run("gcloud",
		"compute",
		"instances",
		"delete",
		name,
		"--zone=europe-north1-c",
		"--quiet",
	)

	return err
}

func run(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	// log.Printf("Command output: %s", cmdOut.String())
	return cmdOut.String(), nil
}
