package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/nicolerobin/kafka-monitor/config"
	"github.com/nicolerobin/kafka-monitor/log"
	"github.com/nicolerobin/kafka-monitor/model/influxdb"
	"github.com/nicolerobin/kafka-monitor/model/kafka"
)

func main() {
	var configFile = flag.String("c", "kafka_monitor.conf", "config file")
	flag.Parse()
	err := config.LoadConfig(*configFile)
	if err != nil {
		panic(err)
	}

	kafka.Init(config.Conf.Kafka.Address)
	influxdb.Init(config.Conf.Influxdb.Address)

	offsets, err := kafka.GetOffsets()
	if err != nil {
		log.Errorf("kafka.GetOffsets() failed, err:%s", err)
		return
	}
	byteResult, err := json.Marshal(offsets)
	if err != nil {
		log.Errorf("json.Marshal() failed, err:%s", err)
		return
	}
	log.Debugf("offsets:%s", string(byteResult))
	for consumerGroup, offset := range offsets {
		for topic, partition := range offset {
			for partitionId, offsetInfo := range partition {
				tags := map[string]string{
					"consumerGroup": consumerGroup,
					"topic":         topic,
					"partitionId":   fmt.Sprintf("%d", partitionId),
				}
				for key, value := range config.Conf.Tags {
					tags[key] = value
				}

				fields := map[string]interface{}{
					"offset":  offsetInfo.Offset,
					"logsize": offsetInfo.Logsize,
					"lag":     offsetInfo.Logsize - offsetInfo.Offset,
				}
				err = influxdb.Report("kafka-monitor", tags, fields)
				if err != nil {
					log.Errorf("influxdb.Report() failed, err:%s", err)
				}
			}
		}
	}
}
