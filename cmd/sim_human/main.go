package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/angel/angel/internal/wire"
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	// ۱. ساخت اتاق
	resp, err := http.Post("http://localhost:8080/rooms", "application/json",
		bytes.NewReader([]byte(`{"game_id":"tic-tac-toe","players":["sim_human", "bot"]}`)))
	if err != nil { fmt.Println("POST_ERR:", err); os.Exit(1) }
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var room struct{ RoomID string `json:"room_id"` }
	json.Unmarshal(data, &room)
	fmt.Println("✅ [HUMAN] Room Created:", room.RoomID)

	// ۲. اتصال به وب‌سوکت
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws?user=sim_human", nil)
	if err != nil { fmt.Println("DIAL_ERR:", err); os.Exit(1) }
	defer conn.Close(websocket.StatusNormalClosure, "")

	// خواندن پاسخ‌های سرور
	go func() {
		for {
			_, data, err := conn.Read(ctx)
			if err != nil { return }
			var sm wire.ServerMsg
			if proto.Unmarshal(data, &sm) == nil {
				if sm.Error != nil {
					fmt.Printf("❌ [SERVER ERROR] %s: %s\n", sm.Error.Code, sm.Error.Error)
				} else if len(sm.Events) > 0 {
					fmt.Printf("📥 [SERVER EVENT] %s\n", sm.Events[0].Type)
				} else {
					fmt.Printf("📥 [SERVER MSG] %s\n", sm.Type)
				}
			}
		}
	}()

	// ۳. ارسال Join
	cmJoin := &wire.ClientMsg{Type: "join", RoomId: room.RoomID}
	bJoin, _ := proto.Marshal(cmJoin)
	conn.Write(ctx, websocket.MessageBinary, bJoin)
	fmt.Println("📤 [HUMAN] Sent JOIN")

	// ۴. صبر برای ورود ربات و شروع بازی
	fmt.Println("⏳ [HUMAN] Waiting 4 seconds for bot & countdown...")
	time.Sleep(4 * time.Second)

	// ۵. ارسال Move (انتخاب خانه 0)
	payload := map[string]any{"cell": 0}
	payloadBytes, _ := json.Marshal(payload)
	cmMove := &wire.ClientMsg{Type: "move", RoomId: room.RoomID, Payload: payloadBytes}
	bMove, _ := proto.Marshal(cmMove)
	fmt.Println("📤 [HUMAN] Sent MOVE (cell: 0)")
	conn.Write(ctx, websocket.MessageBinary, bMove)

	time.Sleep(2 * time.Second)
	fmt.Println("✅ [HUMAN] Test Finished")
}
