package hub

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/coder/websocket"
)

const (
	playTicketLen     = 64
	playGrantTTL      = 30 * time.Second
	playMatchTTL      = 10 * time.Minute
	playJoinTimeout   = 10 * time.Second
	playLobbyGame     = types.GameID("lobby")
)

type playGrant struct {
	User    types.UserID
	Room    types.RoomID
	Expires time.Time
}

type playMatch struct {
	Code    string
	Game    types.GameID
	Public  bool
	Players []types.UserID
	Room    types.RoomID
	Tickets map[types.UserID]string
	Expires time.Time
}

type playLobby struct {
	mu      sync.Mutex
	matches map[string]*playMatch
	grants  map[string]*playGrant
	active  map[string]string
}

type playRequest struct {
	InitData string `json:"init_data"`
	Platform string `json:"platform"`
	Action   string `json:"action"`
	Game     string `json:"game"`
	Code     string `json:"code"`
}

type playResult struct {
	Code   string `json:"code,omitempty"`
	Room   string `json:"room_id,omitempty"`
	User   string `json:"user_id,omitempty"`
	Ticket string `json:"ticket,omitempty"`
	Count  int    `json:"count,omitempty"`
}

func newPlayLobby() *playLobby {
	return &playLobby{
		matches: make(map[string]*playMatch, 64),
		grants:  make(map[string]*playGrant, 128),
		active:  make(map[string]string, 64),
	}
}

func playRandom() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().String()))[:playTicketLen]
	}
	return hex.EncodeToString(b[:])
}

func (p *playLobby) clean(now time.Time) {
	for code, m := range p.matches {
		if now.After(m.Expires) { delete(p.matches, code) }
	}
	for t, g := range p.grants {
		if now.After(g.Expires) { delete(p.grants, t) }
	}
}

func (p *playLobby) authorize(ticket string) (*playGrant, bool) {
	if len(ticket) != playTicketLen { return nil, false }
	p.mu.Lock()
	defer p.mu.Unlock()
	g, ok := p.grants[ticket]
	if !ok { return nil, false }
	if time.Now().After(g.Expires) {
		delete(p.grants, ticket)
		return nil, false
	}
	return g, true
}

func (p *playLobby) begin(h *Hub, uid types.UserID, req playRequest) (*playResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := time.Now()
	p.clean(now)

	gameID := types.GameID(req.Game)
	if gameID == "" { return nil, errBadPlayRequest }

	var rmID types.RoomID
	var code string

	switch req.Action {
	case "create", "":
		room, err := h.OpenRoom(gameID, []types.UserID{uid})
		if err != nil { return nil, err }
		rmID = room.ID()
		code = playRandom()[:8]
		p.matches[code] = &playMatch{
			Code: code, Game: gameID, Public: true, Players: []types.UserID{uid},
			Room: rmID, Tickets: make(map[types.UserID]string, 8), Expires: now.Add(playMatchTTL),
		}

	case "join":
		m, ok := p.matches[req.Code]
		if !ok { return nil, errPlayNoMatch }
		if len(m.Players) >= 8 { return nil, errPlayNoRoom }
		for _, existing := range m.Players {
			if existing == uid {
				rmID = m.Room; code = m.Code
				goto ticket
			}
		}
		m.Players = append(m.Players, uid)
		rmID = m.Room
		code = m.Code
	default:
		return nil, errBadPlayRequest
	}

ticket:
	if rmID == "" { return nil, errPlayNoRoom }
	tk := playRandom()
	p.grants[tk] = &playGrant{User: uid, Room: rmID, Expires: now.Add(playGrantTTL)}
	p.active[tk] = string(uid)

	return &playResult{
		Code: code, Room: string(rmID), User: string(uid), Ticket: tk, Count: len(p.matches),
	}, nil
}

func (h *Hub) servePlayLobby(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{Subprotocols: []string{"angel.play"}})
	if err != nil { return }

	ctx := r.Context()
	for {
		var req playRequest
		if err := wsjsonRead(ctx, conn, &req); err != nil { return }
		res, err := h.play.begin(h, types.UserID(req.InitData), req)
		if err != nil {
			_ = wsjsonWrite(ctx, conn, map[string]string{"error": err.Error()})
			continue
		}
		if err := wsjsonWrite(ctx, conn, res); err != nil { return }
	}
}

func wsjsonRead(ctx context.Context, conn *websocket.Conn, v interface{}) error {
	_, data, err := conn.Read(ctx)
	if err != nil { return err }
	return json.Unmarshal(data, v)
}

func wsjsonWrite(ctx context.Context, conn *websocket.Conn, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil { return err }
	return conn.Write(ctx, websocket.MessageText, data)
}

var _ = logger.Info // 🔴 فقط این خط باقی ماند تا پکیج logger بلااستفاده نماند
