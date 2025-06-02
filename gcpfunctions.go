package a_smarter_email_assistant

import (
	"encoding/json"
	"github.com/jonathanlawhh/a-smarter-email-assistant/AiObject"
	"net/http"
)

// googleFunctionsEntrypoint entry point for Google Functions
func GoogleFunctionsEntrypoint(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusForbidden)
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
