package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"github.com/shuto.sawaki/elmo-project/internal/models"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiAIGenerator は AIGenerator インターフェースを満たす「本物」の実装です。
type GeminiAIGenerator struct {
	model *genai.GenerativeModel
}

// NewGeminiAIGenerator は、Geminiのクライアントを初期化してジェネレータを作成します。
func NewGeminiAIGenerator(ctx context.Context) (*GeminiAIGenerator, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
        // ★ 修正：エラーメッセージを小文字で始める
		return nil, fmt.Errorf("GOOGLE_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
        // ★ 修正：エラーメッセージを小文字で始める
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.Temperature = genai.Ptr(float32(1.0)) // temperatureの設定

	return &GeminiAIGenerator{model: model}, nil
}

// GenerateInitialQuestion は、実際にGemini APIを呼び出して質問を生成します。
func (g *GeminiAIGenerator) SummarizeLogs(ctx context.Context, logs []models.LogEntry) (string, error) {
	// ログエントリのスライスを一つの文字列に結合する
	var logBuilder strings.Builder
	for _, entry := range logs {
		logBuilder.WriteString(entry.Content + "\n")
	}

	prompt := fmt.Sprintf("以下の会議のログを、簡潔で分かりやすい結論として要約してください。重要な決定事項や次のアクションがあれば含めてください。\n\nログ:\n%s", logBuilder.String())

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("gemini api call failed for summarization: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no valid response from gemini api for summarization")
	}

	if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(text), nil
	}

	return "", fmt.Errorf("unexpected response format from gemini api for summarization")
}

// Close はクライアント接続を閉じます。
func (g *GeminiAIGenerator) Close() {
	// genai.Clientには明示的なCloseメソッドがありません。
}

// internal/ai/gemini.go の末尾に追加

// SummarizeLogs は、与えられたログを要約します。
func (g *GeminiAIGenerator) SummarizeLogs(ctx context.Context, logs string) (string, error) {
	prompt := fmt.Sprintf("以下の会議のログを、簡潔で分かりやすい結論として要約してください。重要な決定事項や次のアクションがあれば含めてください。\n\nログ:\n%s", logs)

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("gemini api call failed for summarization: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no valid response from gemini api for summarization")
	}

	if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(text), nil
	}

	return "", fmt.Errorf("unexpected response format from gemini api for summarization")
}