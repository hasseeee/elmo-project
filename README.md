# Elmo Project - 会議室管理システム

## 概要

Elmo Project は、AI を活用した会議室管理システムです。会議の進行、参加者管理、「それな」機能、AI による要約など、モダンな会議体験を提供します。

## 機能

- **会議室管理**: 会議室の作成、編集、ステータス管理
- **参加者管理**: 会議室への参加者追加・削除
- **AI 機能**:
  - 初期質問の自動生成
  - 会議ログの自動要約
- **「それな」機能**: 会議中のリアクション管理
- **結果表示**: 会議結果の詳細表示

## 技術スタック

- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **AI**: Google Gemini API
- **Documentation**: Swagger/OpenAPI

## セットアップ

### 前提条件

- Go 1.25 以上
- PostgreSQL
- Google Gemini API キー

### インストール

1. リポジトリをクローン

```bash
git clone <repository-url>
cd elmo-project
```

2. 依存関係をインストール

```bash
go mod download
```

3. 環境変数を設定

```bash
export GEMINI_API_KEY="your-api-key"
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="elmo-db"
```

4. アプリケーションをビルド

```bash
go build ./cmd/server
```

5. サーバーを起動

```bash
./server
```

## API ドキュメント

### Swagger UI

アプリケーション起動後、以下の URL で Swagger UI にアクセスできます：

```
http://localhost:8080/swagger/index.html
```

### 主要なエンドポイント

#### 会議室管理

- `GET /rooms` - 会議室一覧取得
- `POST /rooms` - 会議室作成
- `GET /rooms/:id` - 会議室詳細取得
- `POST /rooms/:id/start` - 会議開始
- `PUT /rooms/:id/status` - ステータス更新
- `GET /rooms/:id/result` - 会議結果取得
- `POST /rooms/:id/conclusion` - 結論保存
- `POST /rooms/:id/sorena` - 「それな」処理
- `POST /rooms/:id/summary` - 要約作成

#### ユーザー管理

- `POST /users` - ユーザー作成

#### 参加者管理

- `GET /participants` - 参加者一覧取得
- `POST /participants` - 参加者追加

## データベーススキーマ

### 主要テーブル

- `rooms` - 会議室情報
- `users` - ユーザー情報
- `participants` - 参加者情報
- `chat_logs` - チャットログ
- `sorena_counts` - 「それな」カウント

## Docker

### Docker Compose（推奨）

開発環境を簡単に構築するには、Docker Compose を使用します：

```bash
# 開発環境を起動
./scripts/dev.sh start

# 開発環境を停止
./scripts/dev.sh stop

# ログを表示
./scripts/dev.sh logs

# ヘルプを表示
./scripts/dev.sh help
```

### 単体の Docker

単体の Docker を使用してアプリケーションを実行することもできます：

```bash
docker build -t elmo-project .
docker run -p 8080:8080 elmo-project
```

## 開発

### Swagger ドキュメントの更新

コード変更後、Swagger ドキュメントを再生成するには：

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go
```

### テスト

```bash
go test ./...
```

## ライセンス

Apache 2.0

## サポート

問題や質問がある場合は、GitHub の Issues ページでお知らせください。
