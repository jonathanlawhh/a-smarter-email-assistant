package main

import (
	"encoding/json"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/jonathanlawhh/a-smarter-email-assistant/AiObject"
	"github.com/jonathanlawhh/a-smarter-email-assistant/Gemini"
	"github.com/jonathanlawhh/a-smarter-email-assistant/Helper"
	"log"
	"net/http"
	"os"
)

type InputEmail struct {
	EmailDate string   `json:"emailDate"`
	From      string   `json:"emailFrom"`
	To        []string `json:"emailTo"`
	Cc        []string `json:"emailCc"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
}

type RuleSet struct {
	Rule string `json:"rule"`
}

type MyProfile struct {
	Email           []string `json:"myEmail"`
	MyName          []string `json:"myName"`
	MyReplyOpenings []string `json:"myReplyOpenings"`
	MyReplyClosings []string `json:"myReplyClosings"`
}

type InputPayload struct {
	InputEmailData InputEmail `json:"emailPayload"`
	RuleSetData    []RuleSet  `json:"ruleSet"`
	MyProfileData  MyProfile  `json:"myProfile"`
	NoGoWords      []string   `json:"noGoWords"`
}

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
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var inp InputPayload

	err := json.NewDecoder(r.Body).Decode(&inp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	replyResp, err := generateEmailReply(inp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	respDataByte, _ := json.Marshal(replyResp)
	w.Write(respDataByte)
}

func generateEmailReply(inputPayloadData InputPayload) (AiObject.ResponseObj, error) {
	var finalResponseObj AiObject.ResponseObj

	// Encode the sensitive data
	encodedWordMap := &[]Helper.EncodeMap{}
	inputPayloadData.MyProfileData.Email = Helper.MapWordEncoding(inputPayloadData.MyProfileData.Email, encodedWordMap)
	inputPayloadData.InputEmailData.From = Helper.MapWordEncoding([]string{inputPayloadData.InputEmailData.From}, encodedWordMap)[0] // Because from is a string, we pretend to send in an array of string, and take the first value returned
	inputPayloadData.InputEmailData.To = Helper.MapWordEncoding(inputPayloadData.InputEmailData.To, encodedWordMap)
	inputPayloadData.InputEmailData.Cc = Helper.MapWordEncoding(inputPayloadData.InputEmailData.Cc, encodedWordMap)
	Helper.MapWordEncoding(inputPayloadData.NoGoWords, encodedWordMap)

	// Sanitize Data
	inputPayloadData.InputEmailData.Body = Helper.EncodeWordsInString(inputPayloadData.InputEmailData.Body, encodedWordMap)
	for ir, r := range inputPayloadData.RuleSetData {
		inputPayloadData.RuleSetData[ir].Rule = Helper.EncodeWordsInString(r.Rule, encodedWordMap)
	}

	// Create a new input without No Go Words
	combinedInput := InputPayload{
		InputEmailData: inputPayloadData.InputEmailData,
		RuleSetData:    inputPayloadData.RuleSetData,
		MyProfileData:  inputPayloadData.MyProfileData,
	}

	combinedInputJSON, _ := json.Marshal(combinedInput)

	constructPrompt := string(combinedInputJSON)

	log.Println(constructPrompt)

	finalResponseObj, err := Gemini.GenerateResponse(constructPrompt)
	if err != nil {
		return finalResponseObj, err
	}

	finalResponseObj.ReplyBody = Helper.DecodeWordsInString(finalResponseObj.ReplyBody, encodedWordMap)
	//log.Println(finalResponseObj)

	// We deliver the data used so users know what information that need to obscure in the future
	finalResponseObj.LlmPrompt = constructPrompt

	return finalResponseObj, nil
}

func main() {
	// Use PORT environment variable, or default to 8080.
	functions.HTTP("getReply", googleFunctionsEntrypoint)
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
