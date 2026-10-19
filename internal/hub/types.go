package hub
import ("sync";"time";"github.com/angel/angel/internal/bridge";"github.com/angel/angel/internal/core/types";"github.com/angel/angel/internal/registry";"github.com/angel/angel/internal/room";"github.com/coder/websocket")
const ( writeWait = 10 * time.Second; maxMessageSize = 64 << 10 )
type client struct { hub *Hub; conn *websocket.Conn; out chan []byte; id types.UserID; room *room.Room; mu sync.Mutex }
type Hub struct { reg *registry.Registry; br *bridge.Bridge; StatsRecorder room.StatsRecorder; mu sync.Mutex; clients map[*client]struct{}; rooms map[types.RoomID]*room.Room; done chan struct{}; once sync.Once }
func New(reg *registry.Registry) *Hub { return &Hub{reg: reg, clients: make(map[*client]struct{}), rooms: make(map[types.RoomID]*room.Room), done: make(chan struct{})} }
func (h *Hub) SetBridge(br *bridge.Bridge) { h.br = br }
