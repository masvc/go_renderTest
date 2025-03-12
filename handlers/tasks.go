package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-web-app/db"
	"go-web-app/models"
)

// タスク作成時のリクエストボディの構造体
type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

// UpdateTaskRequest タスク更新時のリクエストボディの構造体
type UpdateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      bool      `json:"status"`
}

// ListTasks タスク一覧を取得するハンドラー
func ListTasks(w http.ResponseWriter, r *http.Request) {
	// GETメソッド以外は受け付けない
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// データベースからタスク一覧を取得
	var tasks []models.Task
	result := db.DB.Find(&tasks)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスク一覧の取得に失敗しました"})
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// CreateTask タスクを作成するハンドラー
func CreateTask(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外は受け付けない
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// リクエストボディを読み取る
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "リクエストボディの解析に失敗しました"})
		return
	}

	// タイトルは必須
	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "タイトルは必須です"})
		return
	}

	// タスクを作成
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Status:      false, // 初期状態は未完了
	}

	// データベースに保存
	result := db.DB.Create(&task)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスクの作成に失敗しました"})
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(task)
}

// GetTask 指定されたIDのタスクを取得するハンドラー
func GetTask(w http.ResponseWriter, r *http.Request) {
	// GETメソッド以外は受け付けない
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// URLからIDを取得 (/tasks/1 の場合、1を取得)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "無効なURL形式です"})
		return
	}

	// IDを数値に変換
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "無効なIDです"})
		return
	}

	// データベースからタスクを取得
	var task models.Task
	result := db.DB.First(&task, id)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスクが見つかりません"})
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTask handles PUT requests to update an existing task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// URLからIDを取得
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "無効なURL形式です"})
		return
	}

	// IDを数値に変換
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "無効なIDです"})
		return
	}

	// リクエストボディの解析
	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "リクエストボディの解析に失敗しました"})
		return
	}

	// タスクの取得
	var task models.Task
	if result := db.DB.First(&task, id); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスクが見つかりません"})
		return
	}

	// タスクの更新
	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.DueDate = &req.DueDate

	if result := db.DB.Save(&task); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスクの更新に失敗しました"})
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTask 指定されたIDのタスクを削除するハンドラー
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// DELETEメソッド以外は受け付けない
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// URLからIDを取得
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "無効なURL形式です"})
		return
	}

	// IDを数値に変換
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "無効なIDです"})
		return
	}

	// タスクの存在確認
	var task models.Task
	if result := db.DB.First(&task, id); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスクが見つかりません"})
		return
	}

	// タスクを削除
	if result := db.DB.Delete(&task); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスクの削除に失敗しました"})
		return
	}

	// 204 No Contentを返す
	w.WriteHeader(http.StatusNoContent)
}

// ToggleTask タスクの完了状態を切り替えるハンドラー
func ToggleTask(w http.ResponseWriter, r *http.Request) {
	// PATCHメソッド以外は受け付けない
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// URLからIDを取得
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[3] != "toggle" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "無効なURL形式です"})
		return
	}

	// IDを数値に変換
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "無効なIDです"})
		return
	}

	// タスクの取得
	var task models.Task
	if result := db.DB.First(&task, id); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスクが見つかりません"})
		return
	}

	// 状態を反転
	task.Status = !task.Status

	// データベースを更新
	if result := db.DB.Save(&task); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "タスクの更新に失敗しました"})
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
