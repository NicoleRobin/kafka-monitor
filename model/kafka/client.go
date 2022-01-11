package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/nicolerobin/kafka-monitor/log"
)

var client sarama.Client
var clusterAdmin sarama.ClusterAdmin

func Init(address string) {
	var err error
	client, err = sarama.NewClient([]string{address}, sarama.NewConfig())
	if err != nil {
		panic(err)
	}

	clusterAdmin, err = sarama.NewClusterAdminFromClient(client)
	if err != nil {
		panic(err)
	}
}

func GetLogsize(topic string, partitionId int32) (int64, error) {
	offset, err := client.GetOffset(topic, partitionId, sarama.OffsetNewest)
	if err != nil {
		log.Errorf("client.GetOffset() failed, err:%s", err)
		return 0, err
	}

	return offset, nil
}

type OffsetInfo struct {
	Offset  int64
	Logsize int64
}

func GetOffsets() (map[string]map[string]map[int32]OffsetInfo, error) {
	// topics
	topics, err := clusterAdmin.ListTopics()
	topicNames := make([]string, len(topics))
	for topic, _ := range topics {
		topicNames = append(topicNames, topic)
	}

	topicMetadatas, err := clusterAdmin.DescribeTopics(topicNames)
	if err != nil {
		log.Errorf("clusterAdmin.DescribeTopics() failed, err:%s", err)
		return nil, err
	}

	topicPartitions := make(map[string][]int32)
	for _, topicMetadata := range topicMetadatas {
		partitions := make([]int32, len(topicMetadata.Partitions))
		for _, partition := range topicMetadata.Partitions {
			partitions = append(partitions, partition.ID)
		}
		if topicMetadata.IsInternal {
			continue
		}
		topicPartitions[topicMetadata.Name] = partitions
	}

	// consumer groups
	cgs, err := clusterAdmin.ListConsumerGroups()
	if err != nil {
		log.Errorf("clusterAdmin.ListConsumerGroups() failed, err:%s", err)
		return nil, err
	}

	result := make(map[string]map[string]map[int32]OffsetInfo, len(cgs))
	for cgKey, _ := range cgs {
		cgos, err := clusterAdmin.ListConsumerGroupOffsets(cgKey, topicPartitions)
		if err != nil {
			log.Errorf("clusterAdmin.ListConsumerGroupOffsets() failed, err:%s", err)
			continue
		}
		for topic, blockValue := range cgos.Blocks {
			for partitionId, partition := range blockValue {
				if partition.Err != sarama.ErrNoError {
					log.Errorf("err:%s", partition.Err)
					continue
				}
				if partition.Offset == -1 {
					continue
				}
				logSize, err := GetLogsize(topic, partitionId)
				if err != nil {
					log.Errorf("err:%s", err)
					continue
				}
				log.Debugf("consumerGroup:%s, topic:%s, partitionId:%d, offset:%+v, logSize:%d", cgKey, topic, partitionId, partition.Offset, logSize)
				if _, ok := result[cgKey]; !ok {
					result[cgKey] = make(map[string]map[int32]OffsetInfo, 0)
				}
				if _, ok := result[cgKey][topic]; !ok {
					result[cgKey][topic] = make(map[int32]OffsetInfo, len(cgos.Blocks))
				}
				result[cgKey][topic][partitionId] = OffsetInfo{Offset: partition.Offset, Logsize: logSize}
			}
		}
	}
	return result, nil
}
