package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	filepath "path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Node map[string]Node
}

type Node struct {
	Port       int
	Validators []string
}

func (c *Config) getConf() *Config {
	confFile := flag.String("conf", "../config/single_quorum.yml", "input yml conf file")
	flag.Parse()

	filename, _ := filepath.Abs(*confFile)
	log.Println(filename)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	var c Config
	c.getConf()

	fmt.Println(c)
}
