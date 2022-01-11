package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	err := LoadConfig("kafka_monitor.conf")
	if err != nil {
		t.Fatal(err)
	}
	if Conf.Kafka.Address != "9.131.83.130:9092" {
		t.Fail()
	}
	if Conf.Influxdb.Address != "9.59.11.105:8086" {
		t.Fail()
	}
}
