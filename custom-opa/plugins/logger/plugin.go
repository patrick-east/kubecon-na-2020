package logger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/logs"
	"github.com/open-policy-agent/opa/util"
	"github.com/sirupsen/logrus"
)

const PluginName = "kafka_logger"

type Config struct {
	Host string `json:"host"`
	Topic string `json:"topic"`
}

type KafkaLogger struct {
	manager *plugins.Manager
	mtx     sync.Mutex
	config  Config
	p       *kafka.Producer
}

func (l *KafkaLogger) Start(ctx context.Context) error {

	var err error
	l.p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": l.config.Host})
	if err != nil {
		l.p = nil
		return err
	}

	l.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateOK})
	return nil
}

func (l *KafkaLogger) Stop(ctx context.Context) {
	l.p.Close()

	l.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})
}

func (l *KafkaLogger) Reconfigure(ctx context.Context, config interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.config.Host != config.(Config).Host {
		l.Stop(ctx)
		if err := l.Start(ctx); err != nil {
			l.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr})
		}
	}
	l.config = config.(Config)
}


// Log is called by the decision logger when a record (event) should be emitted. The logs.EventV1 fields
// map 1:1 to those described in https://www.openpolicyagent.org/docs/latest/management/#decision-log-service-api.
func (l *KafkaLogger) Log(ctx context.Context, event logs.EventV1) error {
	if l.p == nil {
		return errors.New("plugin in invalid state")
	}

	bs, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("unable to mashal decision log event: %s", err.Error())
	}

	// Use a delivery channel to ensure the decision event has been received by
	// the broker *before* returning from `Log()` and allowing the evaluation
	// result to be returned to the client.
	deliveryChan := make(chan kafka.Event)

	err = l.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &l.config.Topic, Partition: kafka.PartitionAny},
		Value:          bs,
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("unable to send decision log message: %s", err.Error())
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %s", m.TopicPartition.Error)
	}

	logrus.WithField("plugin", PluginName).Debugf("Delivered message to topic %s [%d] at offset %v",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	return nil
}

type Factory struct{}

func (Factory) New(m *plugins.Manager, config interface{}) plugins.Plugin {

	m.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})

	return &KafkaLogger{
		manager: m,
		config:  config.(Config),
	}
}

func (Factory) Validate(_ *plugins.Manager, config []byte) (interface{}, error) {
	parsedConfig := Config{}
	err := util.Unmarshal(config, &parsedConfig)
	return parsedConfig, err
}