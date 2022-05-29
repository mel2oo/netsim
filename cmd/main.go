package main

import (
	"flag"
	"netsim/internal/config"
	"netsim/pkg/logger"
)

func main() {
	var (
		path string
	)

	flag.StringVar(&path, "config", "../configs/config.toml", "netsim config path")
	flag.Parse()

	cnf, err := config.Load(path)
	if err != nil {
		panic(err)
	}

	if err := logger.Init(cnf.Logger); err != nil {
		panic(err)
	}

	// if err := app.New().Run(); err != nil {
	// 	panic(err)
	// }
}
