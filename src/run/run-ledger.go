package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	filepath "path/filepath"
	"time"

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
		// if nodesMap[name] == nil {
		// 	continue
		// }
		nodes = append(nodes, nodesMap[name])
	}
	return nodes
}

func run(c *Config) {
	threshold := c.Default["threshold"]
	nodes := make(map[string]Node)

	channels := make(ChannelType)
	for _, v := range c.Node {
		nodes[v.Name] = v
		channels[v.Name] = make(ChannelValueType)
	}

	for _, v := range c.Node {
		ledger := NewLedger(v, GetNodeSlice(nodes, v.Validators), threshold, &channels)
		go ledger.Start()
	}
	v11 := Value{Data: "value11"}
	v12 := Value{Data: "value12"}
	vPool1 := []Value{v11, v12}
	msgFrom1 := SCPNomination{Votes: vPool1, NodeName: "n1"}
	// msgFrom2 := SCPNomination{Votes: vPool1, NodeName: "n2"}
	// msgFrom3 := SCPNomination{Votes: vPool1, NodeName: "n3"}
	// msgFrom4 := SCPNomination{Votes: vPool1, NodeName: "n4"}

	go func() {
		// for {
		channels["n0"] <- msgFrom1
		channels["n2"] <- msgFrom1
		channels["n3"] <- msgFrom1
		channels["n4"] <- msgFrom1
		time.Sleep(time.Second * 100)
		// }
	}()

	// go func() {
	// 	for {
	// 		channels["n0"] <- msgFrom2
	// 		time.Sleep(time.Second * 200)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		channels["n0"] <- msgFrom3
	// 		time.Sleep(time.Second * 400)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		channels["n0"] <- msgFrom4
	// 		time.Sleep(time.Second * 800)
	// 	}
	// }()

	var input string
	fmt.Scanln(&input)
}

func main() {
	var c Config
	c.getConf()

	fmt.Println(c)

	run(&c)
}
