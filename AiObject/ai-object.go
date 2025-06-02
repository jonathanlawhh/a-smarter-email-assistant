package AiObject

type ResponseObj struct {
	ReplyBody string   `json:"ReplyBody"`
	Actions   []string `json:"Actions"`
	LlmPrompt string   `json:"LLMPrompt"`
}
