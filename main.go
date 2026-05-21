package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ClaimDetail struct {
	Code      string `json:"code"`
	Series    string `json:"series"`
	Applicant string `json:"applicant"`
	Phone     string `json:"phone"`
	Cellar    string `json:"cellar"`
	Address   string `json:"address"`
}

var claimDB = map[string]ClaimDetail{
	"1": {
		Code:      "BQ-0827",
		Series:    "十摊7春分系列",
		Applicant: "可乐",
		Phone:     "138 **** 5678",
		Cellar:    "四平村古窖藏",
		Address:   "福建省宁德市屏南县",
	},
	"2": {
		Code:      "BQ-0901",
		Series:    "十摊9秋分系列",
		Applicant: "小明",
		Phone:     "139 **** 1234",
		Cellar:    "云岭古窖",
		Address:   "云南省大理州",
	},
	"3": {
		Code:      "BQ-1024",
		Series:    "十摊10冬至系列",
		Applicant: "阿华",
		Phone:     "137 **** 9988",
		Cellar:    "终南山藏",
		Address:   "陕西省西安市长安区",
	},
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	json.NewEncoder(w).Encode(v)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Go! 访问成功 ✓")
	})

	http.HandleFunc("/api/claim/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			writeJSON(w, Resp{Code: 0, Message: "ok"})
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/claim/")
		if id == "" {
			writeJSON(w, Resp{Code: 400, Message: "id required"})
			return
		}
		detail, ok := claimDB[id]
		if !ok {
			writeJSON(w, Resp{Code: 404, Message: "not found"})
			return
		}
		writeJSON(w, Resp{Code: 0, Message: "ok", Data: detail})
	})

	addr := ":8080"
	log.Printf("Server running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
