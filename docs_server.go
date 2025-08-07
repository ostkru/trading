package main

import (
"fmt"
"log"
"net/http"
"os"
"path/filepath"
)

func main() {
// Устанавливаем порт
port := "8090"

// Устанавливаем корневую директорию
rootDir := "/var/www/gogo/go-mod"

// Создаем файловый сервер
fs := http.FileServer(http.Dir(rootDir))

// Настраиваем маршруты
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
tf("Запрос: %s %s", r.Method, r.URL.Path)
(rootDir, r.URL.Path)
otExist(err) {

tf("🚀 Документация запущена на порту %s\n", port)
fmt.Printf("🌐 Доступна по адресу: http://92.53.64.38:%s\n", port)
fmt.Printf("📋 Файлы:\n")
fmt.Printf("   - http://92.53.64.38:%s/docs.html\n", port)
fmt.Printf("   - http://92.53.64.38:%s/redoc-documentation.html\n", port)
fmt.Printf("   - http://92.53.64.38:%s/index.html\n", port)

log.Fatal(http.ListenAndServe(":"+port, nil))
}

func showFileList(w http.ResponseWriter, rootDir string) {
w.Header().Set("Content-Type", "text/html; charset=utf-8")

html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        .file-list { margin: 20px 0; }
        .file-item { margin: 10px 0; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <h1>📚 API Documentation</h1>
    <div class="file-list">
        <div class="file-item">
            <a href="/docs.html">📖 Краткая документация</a>
        </div>
        <div class="file-item">
            <a href="/redoc-documentation.html">📋 Полная документация (Redoc)</a>
        </div>
        <div class="file-item">
            <a href="/index.html">🏠 Главная страница</a>
        </div>
    </div>
    <p><strong>IP адрес сервера:</strong> 92.53.64.38</p>
    <p><strong>API приложение:</strong> <a href="http://92.53.64.38:8095">http://92.53.64.38:8095</a></p>
</body>
</html>`

fmt.Fprint(w, html)
}
