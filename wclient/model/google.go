package model

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/opentdp/wechat-rest/args"
	"github.com/opentdp/wechat-rest/wclient/cache"
)

func GeminiChat(id, msg string) (string, error) {

	opts := []option.ClientOption{
		option.WithAPIKey(args.LLM.GoogleAiKey),
	}
	if args.LLM.GoogleAiUrl != "" {
		opts = append(opts, option.WithEndpoint(args.LLM.GoogleAiUrl))
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, opts...)
	if err != nil {
		return "", err
	}

	defer client.Close()

	// 构造请求参数

	model := client.GenerativeModel(cache.Models[id])

	cs := model.StartChat()
	cs.History = []*genai.Content{}

	for _, msg := range cache.History[id] {
		role := "model"
		if msg.Role == "user" {
			role = "user"
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}

	// 请求模型接口

	resp, err := cs.SendMessage(ctx, genai.Text(msg))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("未得到预期的结果")
	}

	// 更新历史记录

	item1 := cache.HistoryItem{
		Role:    "user",
		Content: msg,
	}

	item2 := cache.HistoryItem{
		Role:    "model",
		Content: fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]),
	}

	cache.History[id] = append(cache.History[id], item1, item2)

	return item2.Content, nil

}