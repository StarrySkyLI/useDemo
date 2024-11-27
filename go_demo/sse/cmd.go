package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// ModelClient 封装了与大模型交互的功能
type ModelClient struct {
	BaseURL string
	APIKey  string
}

// NewModelClient 创建一个新的 ModelClient
func NewModelClient(baseURL, apiKey string) *ModelClient {
	return &ModelClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

// CallModel 向大模型发送请求并返回响应
func (client *ModelClient) CallModel(prompt string) (string, error) {
	// 构建请求体
	data := map[string]interface{}{
		"prompt":     prompt,
		"max_tokens": 100, // 示例设置，可以根据需求调整
	}

	// 转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

	// 发送请求
	req, err := http.NewRequest("POST", client.BaseURL+"/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.APIKey)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call model: %v", err)
	}
	defer resp.Body.Close()

	// 读取并解析响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// 假设模型响应是 JSON 格式，我们可以根据实际响应结构进行处理
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// 假设模型返回的结果是 "choices" -> "text"
	if text, ok := result["choices"].([]interface{})[0].(map[string]interface{})["text"].(string); ok {
		return text, nil
	}

	return "", fmt.Errorf("unexpected response format")
}

// LogModelRequest 记录模型请求和响应日志
func LogModelRequest(db *sql.DB, prompt, model, response string) error {
	stmt, err := db.Prepare("INSERT INTO logs(prompt, model, response, created_at) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(prompt, model, response, time.Now())
	return err
}

// InitDB 初始化数据库
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./logs.db")
	if err != nil {
		return nil, err
	}

	// 创建日志表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			prompt TEXT,
			model TEXT,
			response TEXT,
			created_at DATETIME
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
func main() {
	// 初始化数据库
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 设置路由
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		// 获取 prompt
		prompt := r.URL.Query().Get("prompt")
		if prompt == "" {
			http.Error(w, "Prompt is required", http.StatusBadRequest)
			return
		}

		// 创建模型客户端
		modelClient := NewModelClient("https://api.modelprovider.com", "your-api-key")

		// 调用模型
		response, err := modelClient.CallModel(prompt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get response from model: %v", err), http.StatusInternalServerError)
			return
		}

		// 记录日志
		if err := LogModelRequest(db, prompt, "model_name", response); err != nil {
			http.Error(w, fmt.Sprintf("Failed to log request: %v", err), http.StatusInternalServerError)
			return
		}

		// 流式返回数据
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for _, char := range response {
			fmt.Fprintf(w, "data: %c\n\n", char)
			w.(http.Flusher).Flush()
			time.Sleep(50 * time.Millisecond)
		}

		fmt.Fprintf(w, "data: [END]\n\n")
		w.(http.Flusher).Flush()
	})

	// 启动服务
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
