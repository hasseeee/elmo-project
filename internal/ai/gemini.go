package ai

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/genai"
)

// GeminiAIGenerator は AIGenerator インターフェースを満たす「本物」の実装です。
type GeminiAIGenerator struct {
	client *genai.GenerativeModel
}

// NewGeminiAIGenerator は、Geminiのクライアントを初期化してジェネレータを作成します。
func NewGeminiAIGenerator(ctx context.Context) (*GeminiAIGenerator, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_API_KEY 環境変数が設定されていません")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("Geminiクライアントの作成に失敗しました: %w", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	return &GeminiAIGenerator{client: model}, nil
}

// GenerateInitialQuestion は、実際にGemini APIを呼び出して質問を生成します。
func (g *GeminiAIGenerator) GenerateInitialQuestion(ctx context.Context, title, description string) (string, error) {
	prompt := fmt.Sprintf("ディスカッションルームの魅力的な最初の問いかけだけを生成してください。トピックに関連した、意義のある議論を促すような質問にしてください。余計な前置きや説明は不要です。\n\nタイトル: %s\n説明: %s", title, description)

	resp, err := g.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("Gemini APIの呼び出しに失敗しました: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini APIから有効な応答が得られませんでした")
	}
	
	if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(text), nil
	}

	return "", fmt.Errorf("Gemini APIの応答形式が予期せぬものです")
}

// Close はクライアント接続を閉じます。
func (g *GeminiAIGenerator) Close() {
	g.client.Stop()
}