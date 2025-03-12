package main

import (
	"log"
	"net/http"
	"strings"

	"go-web-app/db"
	"go-web-app/handlers"
	"go-web-app/models"
)

func main() {
	// データベース接続の初期化
	db.InitDB()

	// マイグレーションの実行（テーブルの作成）
	db.DB.AutoMigrate(&models.Task{})

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
