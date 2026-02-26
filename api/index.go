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
		raceName = "R-1グランプリ2026"
	}
	raceName = width.Widen.String(raceName)

	names := r.URL.Query()["n"]
	if len(names) == 0 {
		names = []string{
			"しんや", "今井らいぱち", "渡辺銀次", "ななまがり初瀬", "さすらいラビー中田", "真輝志", "ルシファー吉岡", "九条ジョー", "トンツカタンお抹茶", "マユリカ", "かが屋加賀", "かが屋賀屋", "ザ・マミィ酒井", "ザ・マミィ林", "ザ・マミィ黒田", "ザ・マミィ大島", "ザ・マミィ三島", "ザ・マミィアベ",
		}
	}

	numHorses := len(names)
	if numHorses > 18 {
		numHorses = 18
	}

	// --- 【動的対応】全体の幅を計算 ---
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

	gateColors := []string{"#FFFFFF", "#000000", "#FF0000", "#0000FF", "#FFFF00", "#00FF00", "#FFA500", "#FFC0CB"}
	textColors := []string{"black", "white", "white", "white", "black", "black", "black", "black"}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="540" viewBox="0 0 %d 540" xmlns="http://www.w3.org/2000/svg">`, totalWidth, totalWidth))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="540" fill="white" stroke="black" stroke-width="2"/>`, totalWidth))

	// --- レース名（常に一番右端） ---
	raceNameX := totalWidth - colWidth
	sb.WriteString(`<g>`)
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="0" width="%d" height="540" fill="#333333" stroke="black"/>`, raceNameX, colWidth))

	titleFontSize := 18
	if len([]rune(raceName)) > 15 {
		titleFontSize = 14
	}
	// テキスト中央は raceNameX + 24
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="270" font-size="%d" fill="white" style="writing-mode:vertical-rl;text-anchor:middle;font-weight:bold">%s</text>`, raceNameX+24, titleFontSize, raceName))
	sb.WriteString(`</g>`)

	// --- 馬柱の描画（レース名の左から順に左へ並べる） ---
	currentHorseIdx := 0
	for gIdx := 0; gIdx < 8; gIdx++ {
		for hInGate := 0; hInGate < gateCounts[gIdx]; hInGate++ {
			if currentHorseIdx >= numHorses {
				break
			}

			// posX の計算：右端(totalWidth)から、レース名(48)と馬の数分だけ左へ戻る
			posX := totalWidth - colWidth - ((currentHorseIdx + 1) * colWidth)

			gNum := gIdx + 1
			hNum := currentHorseIdx + 1
			textX := posX + 24

			sb.WriteString(fmt.Sprintf(`<g id="h%d">`, hNum))
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="0" width="%d" height="540" fill="none" stroke="black"/>`, posX, colWidth))
			// 枠色
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="0" width="48" height="30" fill="%s" stroke="black"/>`, posX, gateColors[gIdx]))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="22" font-size="14" text-anchor="middle" font-weight="bold" fill="%s">%d</text>`, textX, textColors[gIdx], gNum))
			// 馬番
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="30" width="48" height="40" fill="#f0f0f0" stroke="black"/>`, posX))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="58" font-size="22" text-anchor="middle" font-weight="bold">%d</text>`, textX, hNum))
			// 馬名
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="90" font-size="22" style="writing-mode:vertical-rl;text-anchor:start;font-weight:bold">%s</text>`, textX, names[currentHorseIdx]))

			// 予想印エリア
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="400" width="48" height="140" fill="#fafafa" stroke="black"/>`, posX))
			for j := 1; j <= 3; j++ {
				y := 400 + (j * 35)
				sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#ccc" stroke-dasharray="2 2" />`, posX, y, posX+48, y))
			}
			sb.WriteString(`</g>`)
			currentHorseIdx++
		}
	}
	sb.WriteString(`</svg>`)

	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprint(w, sb.String())
}
