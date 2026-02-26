package main

import (
	"fmt"
	"net/http"

	handler "github.com/k-p5w/derby-gen/api"
)

func main() {
	port := "3000"

	// 全てのパス（"/"）を api.Handler に流す設定
	// これで http://localhost:3000/ でも http://localhost:3000/?n=... でも反応します
	http.HandleFunc("/", handler.Handler)

	fmt.Printf("Server started at http://localhost:%s\n", port)

	// Sirの仰った「ファイアウォールが出ない設定」で起動
	err := http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		fmt.Printf("起動に失敗しました: %v\n", err)
	}
}
