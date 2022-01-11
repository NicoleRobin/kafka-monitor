package influxdb

import (
	"time"

	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/nicolerobin/kafka-monitor/log"
)

var client influxdb.Client

func Init(address string) {
	var err error
	client, err = influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: address,
	})
	if err != nil {
		panic(err)
	}
}

func Report(measurement string, tags map[string]string, fields map[string]interface{}) error {
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database: "joox_monitoring",
	})
	pt, err := influxdb.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		log.Errorf("client.NewPoint() failed, err:%s", err)
		return err
	}
	bp.AddPoint(pt)
	err = client.Write(bp)
	if err != nil {
		log.Errorf("client.Write() failed, err:%s", err)
		return err
	}
	return nil
}
