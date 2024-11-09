package main

import (
	"github.com/alxrusinov/shorturl/internal/app"
	"github.com/alxrusinov/shorturl/internal/config"
)

func main() {
	config := config.NewConfig()

	config.Init()

	app.Run(config)

}
