package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

const ChartFile = "Chart.yaml"

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

type Chart struct {
	ApiVersion   string   `yaml:"apiVersion"`
	Name         string   `yaml:"name"`
	Version      string   `yaml:"version"`
	KubeVersion  string   `yaml:"kubeVersion"`
	Description  string   `yaml:"description"`
	Type         string   `yaml:"type"`
	Keywords     []string `yaml:"keywords"`
	Home         string   `yaml:"home"`
	Sources      []string `yaml:"sources"`
	Dependencies []*Chart `yaml:"dependencies"`
	Repository   string   `yaml:"repository"`
	Icon         string   `yaml:"icon"`
	AppVersion   string   `yaml:"appVersion"`
	Deprecated   bool     `yaml:"deprecated"`
}

func setOutput(name, value string) {
	// Use GITHUB_OUTPUT file for modern GitHub Actions
	githubOutput := os.Getenv("GITHUB_OUTPUT")
	if githubOutput != "" {
		f, err := os.OpenFile(githubOutput, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("Warning: could not open GITHUB_OUTPUT: %v", err)
			// Fallback to deprecated method
			fmt.Printf("::set-output name=%s::%s\n", name, value)
			return
		}
		defer f.Close()
		fmt.Fprintf(f, "%s=%s\n", name, value)
	} else {
		// Fallback to deprecated method for backward compatibility
		fmt.Printf("::set-output name=%s::%s\n", name, value)
	}
}

func main() {
	chartPath := path.Join(os.Getenv("INPUT_PATH"), ChartFile)
	fmt.Printf("Reading values from %s\n", chartPath)

	dat, readErr := os.ReadFile(chartPath)
	check(readErr)

	chart := Chart{}
	yamlErr := yaml.Unmarshal(dat, &chart)
	check(yamlErr)

	setOutput("apiVersion", chart.ApiVersion)
	setOutput("name", chart.Name)
	setOutput("version", chart.Version)
	setOutput("kubeVersion", chart.KubeVersion)
	setOutput("description", chart.Description)
	setOutput("type", chart.Type)
	setOutput("keywords", strings.Join(chart.Keywords, ","))
	setOutput("home", chart.Home)
	setOutput("sources", strings.Join(chart.Sources, ","))
	setOutput("repository", chart.Repository)
	setOutput("icon", chart.Icon)
	setOutput("appVersion", chart.AppVersion)
	setOutput("deprecated", fmt.Sprintf("%t", chart.Deprecated))

	for _, dep := range chart.Dependencies {
		setOutput(fmt.Sprintf("dependencies_%s_version", dep.Name), dep.Version)
		setOutput(fmt.Sprintf("dependencies_%s_repository", dep.Name), dep.Repository)
	}
}
