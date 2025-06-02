package AiObject

import (
	"encoding/json"
	"github.com/jonathanlawhh/a-smarter-email-assistant/Helper"
	"log"
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

type ResponseObj struct {
	ReplyBody string   `json:"ReplyBody"`
	Actions   []string `json:"Actions"`
	LlmPrompt string   `json:"LLMPrompt"`
}

func GenerateEmailReply(inputPayloadData InputPayload) (ResponseObj, error) {
	var finalResponseObj ResponseObj

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

	finalResponseObj, err := generateGeminiResponse(constructPrompt)
	if err != nil {
		return finalResponseObj, err
	}

	finalResponseObj.ReplyBody = Helper.DecodeWordsInString(finalResponseObj.ReplyBody, encodedWordMap)
	//log.Println(finalResponseObj)

	// We deliver the data used so users know what information that need to obscure in the future
	finalResponseObj.LlmPrompt = constructPrompt

	return finalResponseObj, nil
}
