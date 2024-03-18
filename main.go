package main

import (
	"openidea-idea-social-media-app/app"
	_ "openidea-idea-social-media-app/config"
)

func main() {
	app.StartFiberApp()
}
