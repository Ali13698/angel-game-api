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
	resp, err := http.Post("http://localhost:8080/rooms", "application/json",
		bytes.NewReader([]byte(`{"game_id":"tic-tac-toe","players":["p1", "p2"]}`)))
	if err != nil { os.Exit(1) }
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var room struct{ RoomID string `json:"room_id"` }
	json.Unmarshal(data, &room)
	fmt.Println("✅ [SIM] Room Created:", room.RoomID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn1, _, _ := websocket.Dial(ctx, "ws://localhost:8080/ws?user=p1", nil)
	defer conn1.Close(websocket.StatusNormalClosure, "")

	conn2, _, _ := websocket.Dial(ctx, "ws://localhost:8080/ws?user=p2", nil)
	defer conn2.Close(websocket.StatusNormalClosure, "")

	go func() {
		for {
			_, data, err := conn1.Read(ctx)
			if err != nil { return }
			var sm wire.ServerMsg
			if proto.Unmarshal(data, &sm) == nil {
				if sm.Error != nil {
					fmt.Printf("❌ [P1 ERROR] %s: %s\n", sm.Error.Code, sm.Error.Error)
				} else if len(sm.Events) > 0 {
					fmt.Printf("📥 [P1 EVENT] %s\n", sm.Events[0].Type)
				} else if sm.Type != "" {
					fmt.Printf("📥 [P1 MSG] %s\n", sm.Type)
				}
			}
		}
	}()

	cmJoin := &wire.ClientMsg{Type: "join", RoomId: room.RoomID}
	bJoin, _ := proto.Marshal(cmJoin)
	conn1.Write(ctx, websocket.MessageBinary, bJoin)
	fmt.Println("📤 [P1] Sent JOIN")
	time.Sleep(500 * time.Millisecond)

	conn2.Write(ctx, websocket.MessageBinary, bJoin)
	fmt.Println("📤 [P2] Sent JOIN")

	fmt.Println("⏳ [SIM] Waiting 4 seconds for countdown to finish...")
	time.Sleep(4 * time.Second)

	payload := map[string]any{"cell": 0}
	payloadBytes, _ := json.Marshal(payload)
	cmMove := &wire.ClientMsg{Type: "move", RoomId: room.RoomID, Payload: payloadBytes}
	bMove, _ := proto.Marshal(cmMove)
	fmt.Println("📤 [P1] Sent MOVE (cell: 0)")
	conn1.Write(ctx, websocket.MessageBinary, bMove)

	time.Sleep(2 * time.Second)
	fmt.Println("✅ [SIM] Test Finished")
}
