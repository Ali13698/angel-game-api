package hub

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/profile"
	"github.com/angel/angel/internal/wire"
	"github.com/coder/websocket"
)

type playGrant struct {
	User    types.UserID
	Room    types.RoomID
	Expires time.Time
}

type playMatch struct {
	Code    string
	Game    string
	Public  bool
	Players []types.UserID
	Room    types.RoomID
	Tickets map[types.UserID]string
	Expires time.Time
}

type playLobby struct {
	mu      sync.Mutex
	matches map[string]*playMatch
	grants  map[string]playGrant
	active  map[types.UserID]bool
	slots   chan struct{}
}

type playRequest struct {
	InitData string `json:"initData"`
	Platform string `json:"platform"`
	Action   string `json:"action"`
	Game     string `json:"game_id"`
	Code     string `json:"code"`
}

type playResult struct {
	Code   string `json:"code"`
	Room   string `json:"room_id,omitempty"`
	User   string `json:"user_id"`
	Ticket string `json:"ticket,omitempty"`
	Count  int    `json:"count"`
}

func newPlayLobby() *playLobby {
	return &playLobby{
		matches: make(map[string]*playMatch),
		grants:  make(map[string]playGrant),
		active:  make(map[types.UserID]bool),
		slots:   make(chan struct{}, 128),
	}
}

func playRandom() (string, error) {
	var bytes [32]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[:]), nil
}

func (p *playLobby) clean(now time.Time) {
	for key, match := range p.matches {
		if now.After(match.Expires) {
			delete(p.matches, key)
		}
	}
	for key, grant := range p.grants {
		if now.After(grant.Expires) {
			delete(p.grants, key)
		}
	}
}

func (p *playLobby) authorize(ticket string) (playGrant, bool) {
	if len(ticket) != 64 {
		return playGrant{}, false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clean(time.Now())
	grant, ok := p.grants[ticket]
	return grant, ok
}

func (p *playLobby) begin(h *Hub, uid types.UserID, req playRequest) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clean(time.Now())

	if p.active[uid] {
		return "", errors.New("این حساب از قبل در صف است")
	}

	var match *playMatch

	switch req.Action {
	case "join":
		match = p.matches[req.Code]
		if match == nil || match.Room != "" {
			return "", errors.New("دعوت منقضی شده یا اتاق پر شده است")
		}
		if match.Game != req.Game {
			return "", errors.New("دعوت متعلق به بازی دیگری است")
		}

	case "quick", "create":
		if req.Game == "" || len(req.Game) > 80 {
			return "", errors.New("شناسه بازی نامعتبر است")
		}
		if req.Action == "quick" {
			for _, candidate := range p.matches {
				if candidate.Public &&
					candidate.Game == req.Game &&
					candidate.Room == "" &&
					len(candidate.Players) == 1 {
					match = candidate
					break
				}
			}
		}
		if match == nil {
			if len(p.matches) >= 256 {
				return "", errors.New("ظرفیت صف تکمیل است")
			}
			code, err := playRandom()
			if err != nil {
				return "", errors.New("ساخت دعوت ممکن نشد")
			}
			match = &playMatch{
				Code: code, Game: req.Game,
				Public:  req.Action == "quick",
				Tickets: make(map[types.UserID]string),
				Expires: time.Now().Add(5 * time.Minute),
			}
			p.matches[code] = match
		}
	default:
		return "", errors.New("درخواست نامعتبر است")
	}

	for _, player := range match.Players {
		if player == uid {
			return "", errors.New("برای تست به حساب دوم نیاز است")
		}
	}
	if len(match.Players) >= 2 {
		return "", errors.New("اتاق پر است")
	}

	match.Players = append(match.Players, uid)

	if len(match.Players) == 2 {
		h.mu.Lock()
		full := len(h.rooms) >= 128
		h.mu.Unlock()

		if full {
			match.Players = match.Players[:1]
			return "", errors.New("ظرفیت اتاق‌های این نوبت تست تکمیل است")
		}

		players := append([]types.UserID(nil), match.Players...)
		rm, err := h.OpenRoom(types.GameID(match.Game), players)
		if err != nil {
			match.Players = match.Players[:1]
			return "", errors.New("بازی در میزبان در دسترس نیست")
		}
		match.Room = rm.ID()
		match.Expires = time.Now().Add(30 * time.Minute)
	}

	p.active[uid] = true
	return match.Code, nil
}

func (p *playLobby) status(code string, uid types.UserID) (playResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clean(time.Now())

	match := p.matches[code]
	if match == nil {
		return playResult{}, errors.New("زمان انتظار تمام شد؛ دوباره وارد صف شوید")
	}

	member := false
	for _, player := range match.Players {
		if player == uid {
			member = true
		}
	}
	if !member {
		return playResult{}, errors.New("عضو این صف نیستید")
	}

	result := playResult{
		Code: code, User: string(uid),
		Count: len(match.Players),
	}

	if match.Room != "" {
		ticket := match.Tickets[uid]
		if ticket == "" {
			var err error
			ticket, err = playRandom()
			if err != nil {
				return playResult{}, errors.New("صدور مجوز اتصال ممکن نشد")
			}
			match.Tickets[uid] = ticket
			p.grants[ticket] = playGrant{
				User: uid, Room: match.Room,
				Expires: match.Expires,
			}
		}
		result.Room = string(match.Room)
		result.Ticket = ticket
	}

	return result, nil
}

func (p *playLobby) release(code string, uid types.UserID) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.active, uid)

	match := p.matches[code]
	if match == nil || match.Room != "" {
		return
	}
	for i, player := range match.Players {
		if player == uid {
			match.Players = append(match.Players[:i], match.Players[i+1:]...)
			break
		}
	}
	if len(match.Players) == 0 {
		delete(p.matches, code)
	}
}

func (h *Hub) servePlayLobby(w http.ResponseWriter, r *http.Request) {
	select {
	case h.play.slots <- struct{}{}:
		defer func() { <-h.play.slots }()
	default:
		http.Error(w, "lobby busy", http.StatusServiceUnavailable)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()
	conn.SetReadLimit(32768)

	send := func(message wire.OutboundMsg) error {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		return conn.Write(ctx, websocket.MessageBinary, message.Bytes())
	}
	reject := func(message string) {
		_ = send(wire.OutboundMsg{
			Type: wire.TypeError, Code: wire.CodeRejected, Error: message,
		})
	}

	firstCtx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	kind, data, err := conn.Read(firstCtx)
	cancel()
	if err != nil {
		return
	}
	if kind != websocket.MessageBinary {
		reject("فقط Protobuf باینری پذیرفته می‌شود")
		return
	}

	message, err := wire.NewInboundMsg(data)
	if err != nil || message.Type != "match" {
		reject("درخواست ورود به صف نامعتبر است")
		return
	}

	var req playRequest
	if json.Unmarshal(message.Payload, &req) != nil {
		reject("اطلاعات صف نامعتبر است")
		return
	}

	id, ok := profile.GameIdentity(req.InitData, req.Platform)
	req.InitData = ""
	if !ok {
		reject("ورود معتبر نیست یا منقضی شده؛ مینی‌اپ را از ربات دوباره باز کنید")
		return
	}
	uid := types.UserID(id)

	code, err := h.play.begin(h, uid, req)
	if err != nil {
		reject(err.Error())
		return
	}
	defer h.play.release(code, uid)

	ctx, stop := context.WithTimeout(r.Context(), 5*time.Minute)
	defer stop()

	for {
		result, err := h.play.status(code, uid)
		if err != nil {
			reject(err.Error())
			return
		}

		frameType := "waiting"
		if result.Room != "" {
			frameType = "matched"
		}
		if send(wire.OutboundMsg{Type: frameType, Payload: result}) != nil {
			return
		}
		if frameType == "matched" {
			_ = conn.Close(websocket.StatusNormalClosure, "matched")
			return
		}

		started := time.Now()
		kind, data, err = conn.Read(ctx)
		if err != nil {
			return
		}
		message, err = wire.NewInboundMsg(data)
		if kind != websocket.MessageBinary || err != nil || message.Type != "status" {
			reject("درخواست وضعیت نامعتبر است")
			return
		}
		if time.Since(started) < 400*time.Millisecond {
			reject("درخواست‌ها بیش از حد سریع هستند")
			return
		}
	}
}
