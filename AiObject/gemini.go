package AiObject

import (
	"context"
	"encoding/json"
	"google.golang.org/genai"
	"os"
	"strings"
)

func generateGeminiResponse(constructPrompt string) (ResponseObj, error) {
	var outputData ResponseObj

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return outputData, err
	}

	// List of system instruction. Easier to manage in an array then joined into a string.
	systemInstructionList := []string{
		"You are helping me reply to an email.",
		"I will send you only a JSON payload with the received email, rule list to follow, and my profile.",
		"Based on the payload, generate a reply.",
		"Do not derive name from email.",
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(strings.Join(systemInstructionList, " "), genai.RoleUser),
		ResponseMIMEType:  "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"ReplyBody": {Type: genai.TypeString},
				"Actions": {
					Type:  genai.TypeArray,
					Items: &genai.Schema{Type: genai.TypeString, Enum: []string{"reply_email", "create_task", "to_delete", "do_nothing"}},
				},
			},
			PropertyOrdering: []string{"ReplyBody", "Actions"},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(constructPrompt),
		config,
	)
	if err != nil {
		return outputData, err
	}

	err = json.Unmarshal([]byte(result.Text()), &outputData)
	if err != nil {
		return outputData, err
	}

	// GET metrics
	responseMetrics := struct {
		ModelVersion    string `json:"modelVersion"`
		TotalTokenCount int    `json:"totalTokenCount"`
	}{}
	responseMetricsByte, err := result.MarshalJSON()
	if err != nil {
		return outputData, err
	}
	err = json.Unmarshal(responseMetricsByte, &responseMetrics)
	if err != nil {
		return outputData, err
	}

	return outputData, nil
}
