# シンプルGo Webアプリケーション

## プロジェクト概要
このプロジェクトは、Renderへのデプロイ練習用の基本的なGo言語Webアプリケーションです。
PostgreSQLを使用したCRUD操作の基本と、Renderへのデプロイフローを学ぶことができます。

デモサイト: https://go-rendertest.onrender.com/tasks  
※ 無料プランを使用しているため、非アクティブ時にインスタンスがスピンダウンします。  
　その場合、最初のリクエストに50秒程度の遅延が発生する可能性があります。

## 技術スタック
- 言語: Go
- Webフレームワーク: 標準ライブラリ `net/http`
- データベース: PostgreSQL
- ORM: GORM
- フロントエンド: Bootstrap 5
- テンプレートエンジン: `html/template`
- デプロイ先: Render.com

## 実装状況
### フェーズ1: 基本実装 ✅
- [x] 基本的なCRUD機能の実装
  - [x] タスクの作成 (Create)
  - [x] タスク一覧の取得 (Read)
  - [x] タスクの更新 (Update)
  - [x] タスクの削除 (Delete)
  - [x] タスクの完了状態の切り替え
- [x] PostgreSQLとの連携
- [x] Webインターフェースの実装
  - [x] タスク一覧表示
  - [x] 新規作成フォーム
  - [x] 編集フォーム
  - [x] 削除確認
  - [x] 完了状態の切り替え

### フェーズ2: Renderへのデプロイ ✅
- [x] Renderでのデータベース作成
- [x] アプリケーションのデプロイ設定
- [x] 環境変数の設定
- [x] デプロイの実行と動作確認

## 機能一覧
### エンドポイント
- `GET /` → `/tasks`へリダイレクト
- `GET /health` - ヘルスチェック
- `GET /tasks` - タスク一覧（HTML/JSON）
- `POST /tasks` - タスク作成
- `GET /tasks/:id` - タスク詳細取得
- `PUT /tasks/:id` - タスク更新
- `DELETE /tasks/:id` - タスク削除
- `PATCH /tasks/:id/toggle` - タスク完了状態の切り替え

### Webインターフェース
- タスク一覧表示
  - ID、タイトル、説明、期限、状態を表示
  - 各タスクの編集、削除、状態切り替えが可能
- 新規タスク作成
  - モーダルフォームで作成
  - タイトル（必須）、説明、期限を入力
- タスク編集
  - モーダルフォームで編集
  - タイトル、説明、期限を更新可能

## ローカル開発手順
1. リポジトリのクローン
2. PostgreSQLのインストールと起動
3. データベースの作成:
   ```bash
   createdb go_web_app
   ```
4. 環境変数の設定:
   ```bash
   export DATABASE_URL="postgresql://postgres:password@localhost:5432/go_web_app?sslmode=disable"
   ```
5. 依存パッケージのインストール:
   ```bash
   go mod tidy
   ```
6. アプリケーションの起動:
   ```bash
   go run main.go
   ```
7. ブラウザで http://localhost:8080 にアクセス

## 次のステップ
1. Renderへのデプロイ
2. テストの追加
3. ユーザー認証の実装
4. タスクの検索機能
5. タグ/カテゴリ機能 