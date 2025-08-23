package ai

import (
	"context"
	"github.com/shuto.sawaki/elmo-project/internal/models" // ★ modelsをインポート
)

// AIGenerator は、AIが実行する全ての機能のインターフェースを定義します。
type AIGenerator interface {
	// ★ GenerateInitialQuestion を再度追加
	GenerateInitialQuestion(ctx context.Context, title, description string) (string, error)
	SummarizeLogs(ctx context.Context, logs []models.LogEntry) (string, error)
}