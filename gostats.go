package main

import (
	"encoding/json"
	"fmt"
	goopt "github.com/droundy/goopt"
	"github.com/packetbeat/elastigo/api"
	"github.com/packetbeat/elastigo/core"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type UDPData struct {
	prefix  string
	project string `json: "project"`
	metric  string `json: "metric"`
	value   string `json: "value"`
}

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

func sendToGraphite(message []byte, conn net.UDPConn, graphite net.UDPAddr) {
	if message != nil && len(message) > 0 {
		conn.WriteToUDP([]byte(fmt.Sprintf("%s|g", string(message))), &graphite)
	}
}

func sendToElasticsearch(message []byte) {
	data := createDataStruct(message)
	core.Index(data.prefix, "metric", "1", nil, data)
}

func createDataStruct(message []byte) UDPData {
	valueSplit := strings.Split(string(message), ":")
	keySplit := strings.Split(valueSplit[0], ".")
	return UDPData{project: keySplit[2], prefix: keySplit[0] + "." + keySplit[1], metric: keySplit[3], value: valueSplit[1]}
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

	api.Domain = string(elasticsearch.IP)
	api.Port = strconv.Itoa(elasticsearch.Port)

	graphite, err := getUDPAddressFromConfig("graphite", cfg)
	if err != nil {
		fmt.Printf("Failed to get Server address for Graphite: %v\n", err)
		os.Exit(0)
	}

	graphiteConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		fmt.Printf("Failed to create UDP Connection: %v\n", err)
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
		sendToGraphite(message, *graphiteConn, graphite)
		sendToElasticsearch(message)
	}
}
