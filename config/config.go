package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type KafkaConfig struct {
	Address string `yaml:"address"`
}
type InfluxdbConfig struct {
	Address string `yaml:"address"`
}
type Config struct {
	Kafka    KafkaConfig       `yaml:"kafka"`
	Influxdb InfluxdbConfig    `yaml:"influxdb"`
	Tags     map[string]string `yaml:"tags"`
}

func LoadKafkaConfig() error {
	return nil
}

func LoadInfluxdbConfig() error {
	return nil
}

var Conf *Config

func LoadConfig(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	Conf = &Config{}
	err = yaml.Unmarshal(data, Conf)
	if err != nil {
		return err
	}
	return nil
}
