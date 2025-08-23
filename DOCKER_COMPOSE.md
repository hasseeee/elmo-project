# Docker Compose クイックスタートガイド

## 概要

このガイドでは、Elmo Project を Docker Compose を使用して簡単に起動する方法を説明します。

## 前提条件

- Docker Desktop がインストールされていること
- Docker Compose が利用可能であること

## クイックスタート

### 1. 環境変数の設定

```bash
# 環境変数ファイルをコピー
cp docker-compose.env .env

# .envファイルを編集して、GEMINI_API_KEYを設定
# 例: GEMINI_API_KEY=your-actual-api-key-here
```

### 2. 開発環境の起動

```bash
# 開発環境を起動
./scripts/dev.sh start
```

### 3. アクセス

起動後、以下の URL でアクセスできます：

- **アプリケーション**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **pgAdmin**: http://localhost:5050
- **データベース**: localhost:5432

## 便利なコマンド

### 開発環境管理

```bash
# 開発環境を起動
./scripts/dev.sh start

# 開発環境を停止
./scripts/dev.sh stop

# 開発環境を再起動
./scripts/dev.sh restart

# ログを表示
./scripts/dev.sh logs

# ヘルプを表示
./scripts/dev.sh help
```

### アプリケーション管理

```bash
# アプリケーションをビルド
./scripts/dev.sh build

# テストを実行
./scripts/dev.sh test

# データベースに接続
./scripts/dev.sh db
```

### クリーンアップ

```bash
# コンテナとボリュームを削除
./scripts/dev.sh clean
```

## 手動での Docker Compose 操作

### 基本的な操作

```bash
# サービスを起動
docker-compose up -d

# サービスを停止
docker-compose down

# ログを表示
docker-compose logs -f

# 特定のサービスのログを表示
docker-compose logs -f app
docker-compose logs -f postgres
```

### サービスの再起動

```bash
# 特定のサービスを再起動
docker-compose restart app
docker-compose restart postgres

# 全サービスを再起動
docker-compose restart
```

### ボリュームとネットワークの管理

```bash
# ボリュームを削除
docker-compose down -v

# ネットワークを削除
docker-compose down --remove-orphans

# 全体的なクリーンアップ
docker-compose down -v --remove-orphans
docker system prune -f
```

## トラブルシューティング

### よくある問題

#### 1. ポートが既に使用されている

```bash
# 使用中のポートを確認
lsof -i :8080
lsof -i :5432

# プロセスを終了
kill -9 <PID>
```

#### 2. データベース接続エラー

```bash
# データベースの状態を確認
docker-compose ps postgres

# データベースのログを確認
docker-compose logs postgres

# データベースに直接接続
docker-compose exec postgres psql -U postgres -d elmo-db
```

#### 3. アプリケーションのビルドエラー

```bash
# キャッシュをクリアしてビルド
docker-compose build --no-cache app

# ログを確認
docker-compose logs app
```

### ログの確認

```bash
# 全サービスのログ
docker-compose logs

# 特定のサービスのログ
docker-compose logs app
docker-compose logs postgres

# リアルタイムでログを追跡
docker-compose logs -f
```

## 開発のベストプラクティス

### 1. 環境変数の管理

- `.env`ファイルは Git にコミットしない
- 本番環境では適切なシークレット管理を使用
- 開発環境では`docker-compose.env`をテンプレートとして使用

### 2. データの永続化

- PostgreSQL のデータは`postgres_data`ボリュームに保存
- アプリケーションのコードはホストマシンと同期
- Go のキャッシュは`go_cache`ボリュームに保存

### 3. ネットワーク分離

- 全サービスは`elmo-network`で分離
- 外部からのアクセスは必要なポートのみ公開
- サービス間の通信は内部ネットワークを使用

## 本番環境での使用

本番環境では以下の点に注意してください：

1. **セキュリティ**: 適切なパスワードとシークレット管理
2. **バックアップ**: データベースの定期的なバックアップ
3. **監視**: ログとメトリクスの収集
4. **スケーリング**: 必要に応じたサービスのスケール

## サポート

問題が発生した場合は、以下を確認してください：

1. Docker と Docker Compose のバージョン
2. システムリソース（メモリ、ディスク容量）
3. ファイアウォールの設定
4. ログファイルの内容

詳細なログやエラーメッセージとともに、GitHub の Issues ページで報告してください。
