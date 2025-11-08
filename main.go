package main

import (
	"log"
)

func main() {
	app := NewEmotionalSupportApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
