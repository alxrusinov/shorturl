package main

import (
	"fmt"

	"github.com/alxrusinov/shorturl/internal/app"
	"github.com/alxrusinov/shorturl/internal/config"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	config := config.NewConfig()

	config.Init()

	app.Run(config)

}
