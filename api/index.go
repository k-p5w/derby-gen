package handler

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/text/width"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// --- パラメータ取得 ---
	raceName := r.URL.Query().Get("r")
	if raceName == "" {
		raceName = "ごはんのおとも選手権【GII】"
	}
	raceName = width.Widen.String(raceName)

	// 背景色モードの判定：デフォルトON、bg=0のときだけOFF
	useBgColor := r.URL.Query().Get("bg") != "0"

	names := r.URL.Query()["n"]
	if len(names) == 0 {
		names = []string{
			"生卵", "納豆", "明太子", "鮭フレーク", "海苔", "梅干",
			"塩辛", "キムチ", "高菜", "筋子", "佃煮", "生姜焼",
			"出汁巻", "冷奴", "漬物", "味噌汁", "とろろ", "刻み葱",
		}
	}

	numHorses := len(names)
	if numHorses > 18 {
		numHorses = 18
	}

	// --- 全体の幅を計算 ---
	colWidth := 48
	totalWidth := (numHorses + 1) * colWidth

	// --- 枠割計算 ---
	gateCounts := make([]int, 8)
	for i := 0; i < numHorses; i++ {
		if i < 8 {
			gateCounts[i] = 1
		} else {
			gateCounts[7-((i-8)%8)]++
		}
	}

	// 枠の色定義
	gateColors := []string{"#FFFFFF", "#000000", "#FF0000", "#0000FF", "#FFFF00", "#00FF00", "#FFA500", "#FFC0CB"}
	textColors := []string{"black", "white", "white", "white", "black", "black", "black", "black"}

	// うっすらとした背景色の定義（視認性を確保したパステル調）
	bgColors := []string{
		"#f2f2f2", // 1枠: 白（背景と区別するため極薄いグレー）
		"#e0e0e0", // 2枠: 黒（「黒枠」とわかる程度のグレー）
		"#ffebee", // 3枠: 赤
		"#e3f2fd", // 4枠: 青
		"#fffde7", // 5枠: 黄
		"#f1f8e9", // 6枠: 緑
		"#fff3e0", // 7枠: 橙
		"#fce4ec", // 8枠: 桃
	}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="540" viewBox="0 0 %d 540" xmlns="http://www.w3.org/2000/svg">`, totalWidth, totalWidth))

	// SVG全体の背景（ダークモード対策：画像自体に白背景を持たせる）
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="540" fill="white" />`, totalWidth))

	// --- レース名（一番右端） ---
	raceNameX := totalWidth - colWidth
	sb.WriteString(`<g>`)
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="0" width="%d" height="540" fill="#333333" stroke="black"/>`, raceNameX, colWidth))
	titleFontSize := 18
	if len([]rune(raceName)) > 15 {
		titleFontSize = 14
	}
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="270" font-size="%d" fill="white" style="writing-mode:vertical-rl;text-anchor:middle;font-weight:bold">%s</text>`, raceNameX+24, titleFontSize, raceName))
	sb.WriteString(`</g>`)

	// --- 馬柱の描画（右から左へ並べる） ---
	currentHorseIdx := 0
	for gIdx := 0; gIdx < 8; gIdx++ {
		for hInGate := 0; hInGate < gateCounts[gIdx]; hInGate++ {
			if currentHorseIdx >= numHorses {
				break
			}

			posX := totalWidth - colWidth - ((currentHorseIdx + 1) * colWidth)
			gNum := gIdx + 1
			hNum := currentHorseIdx + 1
			textX := posX + 24

			sb.WriteString(fmt.Sprintf(`<g id="h%d">`, hNum))

			// 1. 馬柱全体の背景色
			fillColor := "white"
			if useBgColor {
				fillColor = bgColors[gIdx]
			}
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="0" width="%d" height="540" fill="%s" stroke="black"/>`, posX, colWidth, fillColor))

			// 2. 枠番（最上部）
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="0" width="48" height="30" fill="%s" stroke="black"/>`, posX, gateColors[gIdx]))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="22" font-size="14" text-anchor="middle" font-weight="bold" fill="%s">%d</text>`, textX, textColors[gIdx], gNum))

			// 3. 馬番（上部）
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="30" width="48" height="40" fill="#f0f0f0" stroke="black"/>`, posX))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="58" font-size="22" text-anchor="middle" font-weight="bold">%d</text>`, textX, hNum))

			// 4. 馬名（中央）
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="90" font-size="22" style="writing-mode:vertical-rl;text-anchor:start;font-weight:bold" fill="black">%s</text>`, textX, names[currentHorseIdx]))

			// 5. 予想印エリア（下部：視認性のため少し白に近い色で固定）
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="400" width="48" height="140" fill="#fafafa" fill-opacity="0.8" stroke="black"/>`, posX))
			for j := 1; j <= 3; j++ {
				y := 400 + (j * 35)
				sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#ccc" stroke-dasharray="2 2" />`, posX, y, posX+48, y))
			}
			sb.WriteString(`</g>`)

			currentHorseIdx++
		}
	}

	// 外枠の縁取り
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="540" fill="none" stroke="black" stroke-width="2"/>`, totalWidth))
	sb.WriteString(`</svg>`)

	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprint(w, sb.String())
}
