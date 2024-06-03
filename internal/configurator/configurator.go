package configurator

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultConfigFile = "copier.yaml"

func formatedKey(s string) string {
	ss := strings.Split(s, "_")
	for i := len(ss) - 1; i >= 0; i-- {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}

type Question struct {
	Title   string   `yaml:"help"`
	Default string   `yaml:"default"`
	Choices []string `yaml:"choices"`
}

func Read(filename string) map[string]Question {
	data := make(map[string]Question)
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		panic(err)
	}

	formatedData := make(map[string]Question, len(data))
	for k, v := range data {
		fk := formatedKey(k)
		formatedData[fk] = v
	}
	return formatedData
}
