package main

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/joho/godotenv/autoload"
	a_smarter_email_assistant "github.com/jonathanlawhh/a-smarter-email-assistant"
	"log"
	"os"
)

func main() {
	// Use PORT environment variable, or default to 8080.
	functions.HTTP("getReply", a_smarter_email_assistant.GoogleFunctionsEntrypoint)
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}
	log.Println(port)
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
