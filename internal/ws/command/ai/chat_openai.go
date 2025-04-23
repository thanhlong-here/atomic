package command

import (
	"context"
	"os"

	"atomic/internal/ws"

	openai "github.com/sashabaranov/go-openai"
)

func init() {
	ws.AutoRegister(ChatOpenAI)
}

func ChatOpenAI(msg ws.WSMessage) map[string]interface{} {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return map[string]interface{}{
			"status": "error",
			"error":  "missing OPENAI_API_KEY",
		}
	}

	client := openai.NewClient(apiKey)

	prompt, _ := msg.Payload["prompt"].(string)
	model, ok := msg.Payload["model"].(string)
	if !ok || model == "" {
		model = openai.GPT3Dot5Turbo // mặc định
	}

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}

	return map[string]interface{}{
		"status":  "ok",
		"message": resp.Choices[0].Message.Content,
	}
}
