package ai

import (
	"context"
	"github.com/shuto.sawaki/elmo-project/internal/models"
)

// AIGenerator は、最初の問いかけを生成する機能のインターフェース（役割）を定義します。
type AIGenerator interface {
	GenerateInitialQuestion(ctx context.Context, title, description string) (string, error)
	SummarizeLogs(ctx context.Context, logs []models.LogEntry) (string, error)
}