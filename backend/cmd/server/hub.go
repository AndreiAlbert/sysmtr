package main

import (
	"encoding/json"
	"log"

	pb "github.com/AndreiAlbert/sysmnt/pb"
)

type StreamHub struct {
	clients    map[*Client]bool
	broadcast  chan *pb.SystemStats
	register   chan *Client
	unregister chan *Client
}

func NewStreamHub() *StreamHub {
	return &StreamHub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *pb.SystemStats),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *StreamHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.conn.Close()
			}
		case stats := <-h.broadcast:
			jsonMsg, err := json.Marshal(map[string]any{
				"cpu_usage":   stats.CpuUsage,
				"ram_usage":   stats.RamUsage,
				"instance_id": stats.InstanceId,
				"created_at":  stats.CreatedAt,
			})
			if err != nil {
				log.Printf("failed to marshal stats: %v", err)
				continue
			}
			for client := range h.clients {
				select {
				case client.send <- jsonMsg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
