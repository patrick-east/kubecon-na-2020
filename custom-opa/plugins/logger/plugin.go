package logger

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/logs"
	"github.com/open-policy-agent/opa/util"
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
}

func (l *KafkaLogger) Start(ctx context.Context) error {
	l.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateOK})
	return nil
}

func (l *KafkaLogger) Stop(ctx context.Context) {
	l.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})
}

func (l *KafkaLogger) Reconfigure(ctx context.Context, config interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.config = config.(Config)
}


// Log is called by the decision logger when a record (event) should be emitted. The logs.EventV1 fields
// map 1:1 to those described in https://www.openpolicyagent.org/docs/latest/management/#decision-log-service-api.
func (l *KafkaLogger) Log(ctx context.Context, event logs.EventV1) error {
	_, err := fmt.Fprintln(os.Stderr, event)
	if err != nil {
		l.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr})
	}
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
	fmt.Println(string(config))
	parsedConfig := Config{}
	err := util.Unmarshal(config, &parsedConfig)
	return parsedConfig, err
}