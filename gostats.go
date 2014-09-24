package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	// "github.com/mattbaird/elastigo"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

var port = goopt.Int([]string{"-p", "--port"}, 8125, "UDP Port to use")
var config = goopt.String([]string{"-c", "--config"}, "./config.json", "Configuration file to use")

func readConfig(filepath string) (map[string]interface{}, error) {
	file, e := ioutil.ReadFile(filepath)
	if e != nil {
		return nil, e
	}
	var cfg map[string]interface{}
	if e = json.Unmarshal(file, &cfg); e != nil {
		return nil, e
	}
	return cfg, nil
}

func getUDPAddressFromConfig(serverType string, cfg map[string]interface{}) (net.UDPAddr, error) {
	key := cfg[serverType].(map[string]interface{})
	parsedIp, _, e := net.ParseCIDR(key["ip"].(string))
	if e != nil {
		return net.UDPAddr{}, e
	}
	return net.UDPAddr{IP: parsedIp, Port: int(key["port"].(float64))}, e
}

func getTCPAddressFromConfig(serverType string, cfg map[string]interface{}) (net.TCPAddr, error) {
	key := cfg[serverType].(map[string]interface{})
	parsedIp, _, e := net.ParseCIDR(key["ip"].(string))
	if e != nil {
		return net.TCPAddr{}, e
	}
	return net.TCPAddr{IP: parsedIp, Port: int(key["port"].(float64))}, e
}

func sendToGraphite(message []byte, graphite net.UDPAddr) {

}

func main() {
	goopt.Description = func() string {
		return "Metric Wrapper for (at first) graphite & elasticsearch."
	}
	goopt.Version = "1.0"
	goopt.Summary = "gostats"
	goopt.Parse(nil)

	readConf, err := readConfig(*config)

	if err != nil {
		fmt.Printf("Failed to parse configuration file: %v\n", err)
		os.Exit(1)
	}

	cfg := readConf["config"].(map[string]interface{})

	elasticsearch, err := getTCPAddressFromConfig("elasticsearch", cfg)
	if err != nil {
		fmt.Printf("Failed to get Server address for Elasticsearch: %v\n", err)
		os.Exit(0)
	}

	graphite, err := getUDPAddressFromConfig("graphite", cfg)
	if err != nil {
		fmt.Printf("Failed to get Server address for Graphite: %v\n", err)
		os.Exit(0)
	}

	fmt.Println(elasticsearch)
	fmt.Println(graphite)

	addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(*port))
	fmt.Println(addr)
	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	for {
		message := make([]byte, 512)
		n, _, err := conn.ReadFromUDP(message)
		log.Printf("Got %d bytes\n", n)
		log.Printf("Data: %s", message)
		if err != nil || n == 0 {
			log.Printf("Error is: %s, bytes are: %d", err, n)
			continue
		}
		sendToGraphite(message, graphite)
	}
}
