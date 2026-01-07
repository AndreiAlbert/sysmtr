package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/AndreiAlbert/sysmnt/pb"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func mockWsConnection(t *testing.T) *websocket.Conn {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader.Upgrade(w, r, nil)
	}))
	t.Cleanup(s.Close)
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Failed to open mock ws connection: %s", err)
	}
	return ws
}

func TestStreamHub_RegisterUnregister(t *testing.T) {
	hub := NewStreamHub()
	go hub.run()
	ws := mockWsConnection(t)
	defer ws.Close()
	client := &Client{
		hub:  hub,
		conn: ws,
		send: make(chan []byte, 10),
	}
	hub.register <- client
	assert.Eventually(t, func() bool {
		_, exists := hub.clients[client]
		return exists
	}, 100*time.Millisecond, 10*time.Millisecond, "Client should be registered")

	hub.unregister <- client
	assert.Eventually(t, func() bool {
		_, exists := hub.clients[client]
		return !exists
	}, 100*time.Millisecond, 10*time.Millisecond, "Client should be unregistered")
}

func TestStreamHub_Broadcast(t *testing.T) {
	hub := NewStreamHub()
	go hub.run()

	client := &Client{
		hub:  hub,
		conn: nil,
		send: make(chan []byte, 128),
	}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	stats := &pb.SystemStats{
		InstanceId: "test-instance-1",
		RamUsage:   4096.0,
		CpuUsage:   55.2,
	}
	hub.broadcast <- stats
	select {
	case msg := <-client.send:
		var received map[string]any
		err := json.Unmarshal(msg, &received)
		assert.NoError(t, err)

		assert.Equal(t, "test-instance-1", received["instanceId"])
		assert.Equal(t, 4096.0, received["ramUsage"])
		assert.Equal(t, 55.2, received["cpuUsage"])
	case <-time.After(1 * time.Second):
		t.Fatalf("Timeout: Hub did not send message to client")
	}
}
