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
		bytes.NewReader([]byte(`{"game_id":"ludo","players":["sim1"]}`)))
	if err != nil {
		fmt.Println("POST_ERR:", err)
		os.Exit(1)
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var room struct {
		RoomID string `json:"room_id"`
	}
	json.Unmarshal(data, &room)
	fmt.Println("ROOM CREATED:", room.RoomID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ۲. اتصال به WebSocket
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws?user=sim1", nil)
	if err != nil {
		fmt.Println("DIAL_ERR:", err)
		os.Exit(1)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// خواندن پاسخ‌ها از سرور
	go func() {
		for {
			typ, data, err := conn.Read(ctx)
			if err != nil {
				fmt.Println("READ_ERR:", err)
				return
			}
			fmt.Println("RECV TYPE:", typ)
			if typ == websocket.MessageBinary {
				var sm wire.ServerMsg
				if err := proto.Unmarshal(data, &sm); err == nil {
					fmt.Printf("SERVER MSG: Type=%s, Payload=%s\n", sm.Type, string(sm.Payload))
				}
			}
		}
	}()

	// تابع کمکی برای ارسال پیام Protobuf
	sendProto := func(msgType string, payload map[string]interface{}) {
		payloadBytes, _ := json.Marshal(payload)
		cm := &wire.ClientMsg{
			Type:    msgType,
			Id:      "1",
			RoomId:  room.RoomID,
			Payload: payloadBytes,
		}
		b, err := proto.Marshal(cm)
		if err != nil {
			fmt.Println("MARSHAL_ERR:", err)
			return
		}
		err = conn.Write(ctx, websocket.MessageBinary, b)
		if err != nil {
			fmt.Println("WRITE_ERR:", err)
			return
		}
		fmt.Printf("SENT PROTOBUF: %s\n", msgType)
	}

	// ۳. ارسال دستورات به صورت Protobuf
	sendProto("join", map[string]interface{}{"room_id": room.RoomID})
	time.Sleep(2 * time.Second)

	sendProto("roll", map[string]interface{}{})
	time.Sleep(2 * time.Second)

	sendProto("roll", map[string]interface{}{})
	time.Sleep(2 * time.Second)

	fmt.Println("DONE")
}
