package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/AndreiAlbert/sysmnt/config"
	pb "github.com/AndreiAlbert/sysmnt/pb"
	"github.com/AndreiAlbert/sysmnt/storage"
	"google.golang.org/grpc"
)

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.LoadConfig()
	hub := NewStreamHub()
	store, err := storage.NewPostgresStore(cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	go hub.run()
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterMonitorServiceServer(s, &server{hub: hub, store: store})
		log.Printf("grpc listening on %s(Waiting for agents...)\n", cfg.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(hub, w, r)
	})

	mux.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		history, err := store.GetHistory(context.Background())
		if err != nil {
			fmt.Println("no mers")
		}
		historyJson, err := json.Marshal(history)
		if err != nil {
			fmt.Println("iar no mers")
		}
		w.Write(historyJson)
	})

	log.Printf("HTTP/Ws listening on :%s (Waiting for Browsers...)", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, enableCors(mux)))
}
