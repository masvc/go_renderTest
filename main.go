package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"go-web-app/db"
	"go-web-app/handlers"
	"go-web-app/models"

	"gorm.io/gorm"
)

func main() {
	// データベース接続の初期化
	db.InitDB()

	// マイグレーションの実行（テーブルの作成）
	db.DB.AutoMigrate(&models.Task{})

	// ダミーデータの挿入
	if err := insertDummyData(db.DB); err != nil {
		log.Printf("Warning: Failed to insert dummy data: %v", err)
	}

	// ルートハンドラーの設定
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ルートパスの場合はタスク一覧ページにリダイレクト
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/tasks", http.StatusMovedPermanently)
			return
		}
		// その他のパスは404
		w.WriteHeader(http.StatusNotFound)
	})

	// ヘルスチェックハンドラーの設定
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy"}`))
	})

	// タスク関連のハンドラーを登録
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.Header.Get("Accept") == "application/json" {
				handlers.ListTasks(w, r)
			} else {
				handlers.HandleTasksPage(w, r)
			}
		case http.MethodPost:
			handlers.CreateTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		// /tasks/{id}/toggle の処理
		if len(parts) == 4 && parts[3] == "toggle" {
			if r.Method == http.MethodPatch {
				handlers.ToggleTask(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// /tasks/{id} の処理
		if len(parts) == 3 {
			// POSTメソッドの場合、_methodパラメータをチェック
			if r.Method == http.MethodPost {
				switch r.FormValue("_method") {
				case "PUT":
					handlers.UpdateTask(w, r)
					return
				case "DELETE":
					handlers.DeleteTask(w, r)
					return
				}
			}

			switch r.Method {
			case http.MethodGet:
				handlers.GetTask(w, r)
			case http.MethodPut:
				handlers.UpdateTask(w, r)
			case http.MethodDelete:
				handlers.DeleteTask(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		// /tasks/ へのアクセスは /tasks にリダイレクト
		if r.URL.Path == "/tasks/" {
			http.Redirect(w, r, "/tasks", http.StatusMovedPermanently)
			return
		}
		// その他のパスは404
		w.WriteHeader(http.StatusNotFound)
	})

	// サーバーの起動
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func insertDummyData(db *gorm.DB) error {
	// テーブルが空の場合のみダミーデータを挿入
	var count int64
	if err := db.Model(&models.Task{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // すでにデータが存在する場合は何もしない
	}

	// 日付の作成
	deploy := time.Date(2025, 3, 20, 0, 0, 0, 0, time.UTC)
	test := time.Date(2025, 3, 25, 0, 0, 0, 0, time.UTC)
	auth := time.Date(2025, 3, 30, 0, 0, 0, 0, time.UTC)
	search := time.Date(2025, 4, 5, 0, 0, 0, 0, time.UTC)
	ui := time.Date(2025, 4, 10, 0, 0, 0, 0, time.UTC)

	dummyTasks := []models.Task{
		{
			Title:       "Renderへのデプロイ",
			Description: "Go WebアプリケーションをRenderにデプロイする",
			DueDate:     &deploy,
			Status:      true,
		},
		{
			Title:       "テストコードの作成",
			Description: "アプリケーションの単体テストとE2Eテストを実装する",
			DueDate:     &test,
			Status:      false,
		},
		{
			Title:       "ユーザー認証の実装",
			Description: "JWTを使用したユーザー認証システムを追加する",
			DueDate:     &auth,
			Status:      false,
		},
		{
			Title:       "検索機能の追加",
			Description: "タスクのタイトルと説明で検索できる機能を実装する",
			DueDate:     &search,
			Status:      false,
		},
		{
			Title:       "UIの改善",
			Description: "モバイル対応とダークモードの実装",
			DueDate:     &ui,
			Status:      false,
		},
	}

	return db.Create(&dummyTasks).Error
}
