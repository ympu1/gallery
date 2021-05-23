package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RootPath  string   `yaml:"RootPath"`
	Port      string   `yaml:"Port"`
	ThumbSize int      `yaml:"ThumbSize"`
	ImageExt  []string `yaml:"ImageExt"`
	VideoExt  []string `yaml:"VideoExt"`
}

func (config *config) fillFromYML(ymlFileName string) error {
	content, err := ioutil.ReadFile(ymlFileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}