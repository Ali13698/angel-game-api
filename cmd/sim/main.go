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

    "github.com/coder/websocket"
)

func main() {
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
    fmt.Println("ROOM:", room.RoomID)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws?user=sim1", nil)
    if err != nil {
        fmt.Println("DIAL_ERR:", err)
        os.Exit(1)
    }
    defer conn.Close(websocket.StatusNormalClosure, "")

    go func() {
        for {
            _, data, err := conn.Read(ctx)
            if err != nil {
                fmt.Println("READ_ERR:", err)
                return
            }
            fmt.Println("RECV:", string(data))
        }
    }()

    send := func(v map[string]interface{}) {
        b, _ := json.Marshal(v)
        conn.Write(ctx, websocket.MessageText, b)
        fmt.Println("SENT:", string(b))
    }

    send(map[string]interface{}{"type": "join", "room_id": room.RoomID, "id": "1"})
    time.Sleep(4 * time.Second)
    send(map[string]interface{}{"type": "roll", "id": "2"})
    time.Sleep(2 * time.Second)
    send(map[string]interface{}{"type": "roll", "id": "3"})
    time.Sleep(3 * time.Second)
    fmt.Println("DONE")
}
