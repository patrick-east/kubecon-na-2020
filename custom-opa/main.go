package main

import (
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/cmd"
	"github.com/patrick-east/kubecon-na-2020/custom-opa/builtins"
)

func main() {
	builtins.Register()

	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
