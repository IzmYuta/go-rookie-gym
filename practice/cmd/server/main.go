package main

import (
	"fmt"
	"log"
	"io"
	"encoding/json"
	"net/http"
)


func main() {
	// HandleFuncでエンドポイントと処理を紐付ける
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/hello", HelloServer)

	// ListenAndServeでポート開通
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Payload struct {
	Name string `json:"name"`
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	var p Payload
	// io.ReadAllでリクエストボディを読み込む
	body, _ := io.ReadAll(req.Body)
	// json.UnmarshalでJSONを構造体に変換
	// 変換された構造体はpに格納される
	if err := json.Unmarshal(body, &p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := p.Name
	p.Name = fmt.Sprintf("Hello, %s", name)
	// json.Marshalで構造体をJSONに変換
	j,_ := json.Marshal(&p)
	// JSONをレスポンスとして返す
	w.Write(j)

}