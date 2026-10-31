package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/angel/angel/internal/wire"
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	fmt.Println("=== TTT Simulator Started ===")

	// 1. ساخت اتاق
	resp, err := http.Post("http://localhost:8080/rooms", "application/json",
		bytes.NewReader([]byte(`{"game_id":"tic-tac-toe","players":["sim1", "bot"]}`)))
	if err != nil {
		log.Fatalf("POST_ERR: %v", err)
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var room struct{ RoomID string `json:"room_id"` }
	json.Unmarshal(data, &room)
	fmt.Printf("ROOM CREATED: %s\n", room.RoomID)
	if room.RoomID == "" { log.Fatal("Failed to create room. Response: ", string(data)) }

	// 2. اتصال به وب‌سوکت
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws?user=sim1", nil)
	if err != nil {
		log.Fatalf("WS_DIAL_ERR: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")
	fmt.Println("WS CONNECTED")

	// خواندن پیام‌ها در پس‌زمینه
	go func() {
		for {
			_, data, err := conn.Read(ctx)
			if err != nil { return }
			var sm wire.ServerMsg
			if err := proto.Unmarshal(data, &sm); err == nil {
				evType := "none"
				if len(sm.Events) > 0 { evType = sm.Events[0].Type }
				fmt.Printf("RECV -> Type: %s, Event: %s, Err: %v\n", sm.Type, evType, sm.Error)
			}
		}
	}()

	// 3. ارسال دستور Join
	sendProto(ctx, conn, "join", room.RoomID, map[string]any{})
	time.Sleep(1 * time.Second)

	// 4. ارسال دستور Move (انتخاب خانه 0)
	sendProto(ctx, conn, "move", room.RoomID, map[string]any{"cell": 0})
	time.Sleep(2 * time.Second)

	fmt.Println("=== SIMULATOR DONE ===")
}

func sendProto(ctx context.Context, conn *websocket.Conn, msgType, roomID string, payload map[string]any) {
	payloadBytes, _ := json.Marshal(payload)
	cm := &wire.ClientMsg{
		Type:    msgType,
		Id:      "1",
		RoomId:  roomID,
		Payload: payloadBytes,
	}
	b, _ := proto.Marshal(cm)

	fmt.Printf("SENT -> %s\n", msgType)
	err := conn.Write(ctx, websocket.MessageBinary, b)
	if err != nil { log.Printf("WRITE_ERR: %v\n", err) }
}
