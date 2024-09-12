package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Tempo struct {
	Image         string   `yaml:"image"`
	ContainerName string   `yaml:"container_name"`
	Command       []string `yaml:"command"`
	Volumes       []string `yaml:"volumes"`
	Ports         []string `yaml:"ports"`
	Networks      []string `yaml:"networks"`
}

type DockerCompose struct {
	Tempo Tempo `json:"tempo"`
}

func main() {
	var dockerCom = DockerCompose{
		Tempo: Tempo{
			Image:         "grafana/tempo:3.24",
			ContainerName: "dev-environment-tempo-1",
			Command:       []string{"-config.file=/etc/tempo.yaml"},
			Volumes:       []string{"./configs/tempo/tempo.yaml:/etc/tempo.yaml", "./configs/tempo/data:/tmp/tempo"},
			Ports:         []string{"14268:14268", "3200:3200", "9095:9095", "4317:4317", "4318:4318", "9411:9411"},
			Networks:      []string{"dev-environment"},
		},
	}

	d, err := yaml.Marshal(&dockerCom)
	if err != nil {
		fmt.Println("error:", err)
	}

	file, err := os.Create("tempo.yaml")
	if err != nil {
		return
	}

	_, err = file.Write(d)
	if err != nil {
		return
	}

	file.Close()
}
