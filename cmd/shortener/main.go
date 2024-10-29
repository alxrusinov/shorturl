package main

import (
	"github.com/alxrusinov/shorturl/internal/app"
	"github.com/alxrusinov/shorturl/internal/config"
)

var appConfig *config.Config = config.NewConfig()

func init() {
	appConfig.Init()
}

func main() {
	appConfig.Parse()

	app.Run(appConfig)
}
