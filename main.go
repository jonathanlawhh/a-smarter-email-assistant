package main

import (
	"encoding/json"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/jonathanlawhh/a-smarter-email-assistant/AiObject"
	"log"
	"net/http"
	"os"
)

// googleFunctionsEntrypoint entry point for Google Functions
func googleFunctionsEntrypoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var inp AiObject.InputPayload

	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	replyResp, err := AiObject.GenerateEmailReply(inp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	respDataByte, _ := json.Marshal(replyResp)
	w.Write(respDataByte)
}

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if os.Getenv("PRODUCTION") == "true" {
		log.Print("starting server...")
		http.HandleFunc("/", googleFunctionsEntrypoint)

		// Start HTTP server.
		log.Printf("listening on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	} else {
		// Use PORT environment variable, or default to 8080.
		functions.HTTP("getReply", googleFunctionsEntrypoint)

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
}
