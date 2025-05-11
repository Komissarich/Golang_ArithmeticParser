package main

import (
	"calc/pkg/config"
	"calc/server/agent"
	orchestrator "calc/server/orchestrator"
	_ "net/http/pprof"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		panic("config is missing")
	}

	app := orchestrator.New(*cfg)
	agent := agent.New(cfg.Server_Port, cfg.Agent)
	go agent.Work()
	go app.TaskCreator()
	app.RunServer()

}
