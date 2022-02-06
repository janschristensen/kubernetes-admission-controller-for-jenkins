package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type imageConf struct {
	Default string
	Allow   struct {
		Empty      string
		Registries []string
	}
}

func readConf(filename string) (*imageConf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &imageConf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
