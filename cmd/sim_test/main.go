package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Creating Tic-Tac-Toe room...")
	resp, err := http.Post("http://localhost:8080/rooms", "application/json",
		bytes.NewReader([]byte(`{"game_id":"tic-tac-toe","players":["sim1", "bot"]}`)))
	if err != nil { panic(err) }
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println("Response:", string(data))
	
	var room struct{ RoomID string `json:"room_id"` }
	json.Unmarshal(data, &room)
	
	time.Sleep(5 * time.Second) // صبر برای ورود ربات
	fmt.Println("Room should be active now. Check server logs!")
}
