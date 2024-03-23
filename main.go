package main

import (
	"flag"
	"openidea-idea-social-media-app/app"
	_ "openidea-idea-social-media-app/config"
)

func main() {
	var port string
	var prefork bool
	flag.StringVar(&port, "port", "8080", "application port")
	flag.BoolVar(&prefork, "prefork", false, "enable prefork")
	flag.Parse()

	app.StartFiberApp(port, prefork)
}
