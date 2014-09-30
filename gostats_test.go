package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestGetUDPAddressFromConfig(t *testing.T) {
	cfg := map[string]interface{}{}
	json.Unmarshal([]byte(`
	{
		"udpserver": {
			"ip": "127.0.0.1/24",
			"port": 1234
		}	
	}	
	`), &cfg)
	result, err := getUDPAddressFromConfig("udpserver", cfg)
	expectedIP, _, _ := net.ParseCIDR("127.0.0.1/24")
	assert.Equal(t, net.UDPAddr{IP: expectedIP, Port: 1234}, result)
	assert.Nil(t, err)
}

func TestGetElasticSearchConnection(t *testing.T) {
	cfg := map[string]interface{}{}
	json.Unmarshal([]byte(`
	{
		"elasticsearch": {
			"hostname": "stats.p.local.netconomy.net",
			"port": 1234
		}	
	}	
	`), &cfg)
	conn := getElasticSearchConnection("elasticsearch", cfg)
	assert.Equal(t, conn.Host, "stats.p.local.netconomy.net")
	assert.Equal(t, conn.Port, "1234")
}

func TestCreateDataStruct(t *testing.T) {
	dataStruct := createDataStruct([]byte("sonar.metrics.xlmsp.coverage:15"))
	assert.Equal(t, dataStruct.prefix, "sonar.metrics")
	assert.Equal(t, dataStruct.project, "xlmsp")
	assert.Equal(t, dataStruct.metric, "coverage")
	assert.Equal(t, dataStruct.value, "15")
}
