package main

import (
	"encoding/json"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

type KubectlVersionJSON struct {
	ServerVersion struct {
		Major      string `json:"major"`
		Minor      string `json:"minor"`
		GitVersion string `json:"gitVersion"`
	} `json:"serverVersion"`
}

func getK8sVersion(contextName string) (string, string, string) {
	out, err := exec.Command(config.kubeCtlLocation, "--context", contextName, "version", "--output", "json").Output()
	if err != nil {
		logrus.Errorf("Could not get Kubernetes version for %s", contextName)
		return "", "", ""
	}

	outParsed := KubectlVersionJSON{}
	json.Unmarshal(out, &outParsed)

	return outParsed.ServerVersion.GitVersion, outParsed.ServerVersion.Major, outParsed.ServerVersion.Minor
}
