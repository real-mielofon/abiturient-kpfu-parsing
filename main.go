package main

import (
	"os"
	"time"

	"github.com/real-mielofon/abiturient-kpfu-parsing/app"
	"github.com/real-mielofon/abiturient-kpfu-parsing/config"
)

const (
	fileConfig   = "./data/subscribe.txt"
	periodUpdate = 30 * time.Minute
)

func main() {

	env := os.Getenv("TGBOT_KEY")

	cfg := new(config.Config)
	cfg.ReadConfig(fileConfig)

	a := app.New(cfg, env, periodUpdate)
	a.Run()
}
