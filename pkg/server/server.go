package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type Server struct {
	HTTPPort   int
	GRPCPort   int
	Mux        *http.ServeMux
	uiFS       fs.FS
	wsClients  map[*websocket.Conn]bool
	wsMu       sync.Mutex
	grpcServer *grpc.Server
	httpServer *http.Server
	upgrader   websocket.Upgrader
}

func New(httpPort, grpcPort int, uiAssets fs.FS) *Server {
	return &Server{
		HTTPPort:  httpPort,
		GRPCPort:  grpcPort,
		Mux:       http.NewServeMux(),
		uiFS:      uiAssets,
		wsClients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.Mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	s.Mux.HandleFunc("/ws", s.handleWebSocket)

	if s.uiFS != nil {
		fileServer := http.FileServer(http.FS(s.uiFS))
		s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/ws" {
				http.NotFound(w, r)
				return
			}

			path := r.URL.Path
			if path == "/" {
				path = "/index.html"
			}

			if _, err := fs.Stat(s.uiFS, strings.TrimPrefix(path, "/")); err != nil {
				r.URL.Path = "/"
			}

			fileServer.ServeHTTP(w, r)
		})
	}

	addr := fmt.Sprintf(":%d", s.HTTPPort)
	s.httpServer = &http.Server{Addr: addr, Handler: s.corsMiddleware(s.Mux)}

	go func() {
		<-ctx.Done()
		s.httpServer.Shutdown(context.Background())
	}()

	log.Printf("HTTP server → http://localhost:%d", s.HTTPPort)
	return s.httpServer.ListenAndServe()
}

func (s *Server) StartGRPC(register func(*grpc.Server)) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.GRPCPort))
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	register(s.grpcServer)

	log.Printf("gRPC server → localhost:%d", s.GRPCPort)
	return s.grpcServer.Serve(lis)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.wsMu.Lock()
	s.wsClients[conn] = true
	s.wsMu.Unlock()

	defer func() {
		s.wsMu.Lock()
		delete(s.wsClients, conn)
		s.wsMu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (s *Server) Broadcast(data interface{}) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}

	s.wsMu.Lock()
	defer s.wsMu.Unlock()

	for client := range s.wsClients {
		if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
			client.Close()
			delete(s.wsClients, client)
		}
	}
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
