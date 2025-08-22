package ai

import "context"

// AIGenerator は、最初の問いかけを生成する機能のインターフェース（役割）を定義します。
type AIGenerator interface {
	GenerateInitialQuestion(ctx context.Context, title, description string) (string, error)
}