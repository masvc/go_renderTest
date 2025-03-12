package handlers

import (
	"html/template"
	"net/http"

	"go-web-app/db"
	"go-web-app/models"
)

var templates = template.Must(template.ParseFiles(
	"templates/layouts/base.html",
	"templates/tasks/index.html",
))

// PageData テンプレートに渡すデータの構造体
type PageData struct {
	Title string
	Tasks []models.Task
}

// HandleTasksPage タスク一覧ページのハンドラー
func HandleTasksPage(w http.ResponseWriter, r *http.Request) {
	// タスク一覧の取得
	var tasks []models.Task
	result := db.DB.Find(&tasks)
	if result.Error != nil {
		http.Error(w, "タスクの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	// テンプレートにデータを渡して描画
	data := PageData{
		Title: "タスク一覧",
		Tasks: tasks,
	}

	if err := templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
