package main

import (
	"os"

	_ "github.com/ballot/internals/app"
	"github.com/ballot/internals/routes"
)

func main() {
	// setup cors
	origins := []string{
		os.Getenv("MINIAPP_URL"),
		os.Getenv("FRONTEND_URL"),
		os.Getenv("MINIAPP_TG_URL"),
		os.Getenv("LOCAL_FRONTEND_URL"),
	}

	routes.Cors(origins...)

	// api routes
	routes.RegisteredRoutes()
	routes.Run(":8003")
}
