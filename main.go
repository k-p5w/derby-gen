package main

import (
	"fmt"
	"net/http"

	handler "github.com/k-p5w/derby-gen/api"
)

func main() {
	port := "3000"

	// 1. 画像生成用（/api/generate にアクセスしたとき）
	http.HandleFunc("/api/generate", handler.Handler)

	// 2. 入力画面用（それ以外のパス、つまりトップページにアクセスしたとき）
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// index.html をブラウザに返す
		http.ServeFile(w, r, "index.html")
	})

	fmt.Printf("Server started at http://localhost:%s\n", port)

	// localhost 指定で起動（セキュリティ警告が出にくい設定）
	err := http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		fmt.Printf("起動に失敗しました: %v\n", err)
	}
}
