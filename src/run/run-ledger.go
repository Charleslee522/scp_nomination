package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	filepath "path/filepath"

	. "github.com/Charleslee522/scp_nomination/src/common"
	. "github.com/Charleslee522/scp_nomination/src/ledger"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Default map[string]int
	Node    map[string]Node
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

func GetNodeSlice(nodesMap map[string]Node, validators []string) []Node {
	nodes := []Node{}
	for _, name := range validators {
		nodes = append(nodes, nodesMap[name])
	}
	return nodes
}

func run(c *Config) {
	threshold := c.Default["threshold"]
	nodes := make(map[string]Node)

	for _, v := range c.Node {
		nodes[v.Name] = v
	}

	ledgers := []Ledger{}
	for _, v := range c.Node {
		ledger := NewLedger(v, GetNodeSlice(nodes, v.Validators), threshold)
		ledger.Consensus.InsertValues(v.Messages)
		ledgers = append(ledgers, *ledger)
		go func(ledger *Ledger) {
			ledger.Start()
		}(ledger)
	}

	var input string
	fmt.Scanln(&input)

	for _, ledger := range ledgers {
		fmt.Print(ledger.Node.Name, " has values -> ")
		fmt.Println(ledger.Consensus.GetConfirmedValues())
	}
}

func main() {
	var c Config
	c.getConf()

	run(&c)
}
