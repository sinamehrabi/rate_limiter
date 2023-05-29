package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Path string `yaml:"path"`
	Unit string `yaml:"unit"`
	Rpu  int    `yaml:"rpu"`
}

type conf struct {
	Rules []Rule
}

func getYamlFile(filename string) (map[string]Rule, error) {

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &conf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	rules := make(map[string]Rule)
	for _, v := range c.Rules {
		rules[v.Path] = v
	}

	return rules, err
}
