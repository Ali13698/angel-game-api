package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coder/websocket"
)

type msg struct {
	Type   string         `json:"type"`
	RoomID string         `json:"room_id,omitempty"`
	ID     string         `json:"id,omitempty"`
	Rev    uint64         `json:"rev,omitempty"`
	Payload map[string]any `json:"payload,omitempty"`
}

func main() {
	roomID := os.Args[1]
	userID := os.Args[2]
	mode := "human" // human | bot
	if len(os.Args) > 3 { mode = os.Args[3] }

	ctx := context.Background()
	url := "ws://localhost:8080/ws?user=" + userID
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil { log.Fatal(err) }
	defer c.CloseNow()
	c.SetReadLimit(1 << 20)

	// خواندن پس‌زمینه
	go func() {
		for {
			_, data, err := c.Read(ctx)
			if err != nil { return }
			var out map[string]any
			_ = json.Unmarshal(data, &out)
			t, _ := out["type"].(string)
			fmt.Println("["+userID+"]", t, "|", string(data)[:min(200, len(data))])
		}
	}()

	// join
	send(ctx, c, msg{Type: "join", RoomID: roomID, ID: "1"})
	time.Sleep(4 * time.Second)

	if mode == "bot" {
		// خودکار بازی
		for i := 0; i < 50; i++ {
			send(ctx, c, msg{Type: "roll", ID: fmt.Sprintf("r%d", i)})
			time.Sleep(2 * time.Second)
		}
	} else {
		// انسانی: چند حرکت با تأخیر
		send(ctx, c, msg{Type: "chat", Payload: map[string]any{"text": "سلام از " + userID}})
		time.Sleep(500 * time.Millisecond)
		send(ctx, c, msg{Type: "roll", ID: "r1"})
		time.Sleep(2 * time.Second)
		send(ctx, c, msg{Type: "move", Payload: map[string]any{"token": 0}})
		time.Sleep(3 * time.Second)
		send(ctx, c, msg{Type: "roll", ID: "r2"})
		time.Sleep(2 * time.Second)
		send(ctx, c, msg{Type: "move", Payload: map[string]any{"token": 1}})
		time.Sleep(3 * time.Second)
	}
	fmt.Println("["+userID+"] DONE")
}

func send(ctx context.Context, c *websocket.Conn, m msg) {
	b, _ := json.Marshal(m)
	wctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	_ = c.Write(wctx, websocket.MessageText, b)
	cancel()
	fmt.Println("["+m.ID+"] SENT:", string(b))
}

func min(a, b int) int { if a < b { return a }; return b }
