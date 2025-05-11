package main

import (
	"calc/pkg/config"
	"calc/server/agent"
	orchestrator "calc/server/orchestrator"
<<<<<<< HEAD
	_ "net/http/pprof"
=======
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
)

func main() {

<<<<<<< HEAD
	cfg, err := config.New()
	if err != nil {
		panic("config is missing")
	}

=======
	//calculator.Calc("2+2*3")
	cfg, err := config.New()
	if err != nil {
		//		panic("config is missing")
	}
>>>>>>> c1a028191862e07aa216c4e0bb0d68ac4c4fa868
	app := orchestrator.New(*cfg)
	agent := agent.New(cfg.Server_Port, cfg.Agent)
	go agent.Work()
	go app.TaskCreator()
	app.RunServer()

}
