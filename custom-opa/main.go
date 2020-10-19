package main

import (
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/cmd"
	"github.com/patrick-east/kubecon-na-2020/custom-opa/builtins"
	"github.com/patrick-east/kubecon-na-2020/custom-opa/plugins"
)

func main() {
	builtins.Register()
	plugins.Register()

	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
