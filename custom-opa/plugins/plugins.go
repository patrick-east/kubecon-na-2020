package plugins

import (
	"github.com/open-policy-agent/opa/runtime"
	"github.com/patrick-east/kubecon-na-2020/custom-opa/plugins/api"
	"github.com/patrick-east/kubecon-na-2020/custom-opa/plugins/logger"
)

func Register() {
	runtime.RegisterPlugin(logger.PluginName, logger.Factory{})
	runtime.RegisterPlugin(api.PluginName, api.Factory{})
}
