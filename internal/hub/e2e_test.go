package hub_test

import (
	"context"
	"testing"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/wasm"
	"github.com/angel/angel/internal/room"
)

type mockTransport struct{}
func (m *mockTransport) Broadcast(rmID types.RoomID, frame []byte) {}
func (m *mockTransport) SendTo(rmID types.RoomID, uid types.UserID, frame []byte) {}
func (m *mockTransport) CloseSlowClient(rmID types.RoomID, uid types.UserID) {}
func (m *mockTransport) SignalBackpressure(rmID types.RoomID, level int) {}

func TestGameConnection_E2E(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reg := registry.New(ctx)
	factory := wasm.NewFactory("tic-tac-toe")
	if err := reg.Register(factory); err != nil {
		t.Fatalf("Failed to register factory: %v", err)
	}
	
	transport := &mockTransport{}
	rm, err := reg.Open(registry.Spec{
		GameID: "tic-tac-toe",
		Players: []types.UserID{"p1", "p2"},
		Transport: transport,
	})
	if err != nil {
		t.Fatalf("Failed to open room: %v", err)
	}

	if rm == nil {
		t.Fatal("Expected room instance, got nil")
	}
}
