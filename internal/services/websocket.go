package services

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type WSNotification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type WebSocketService struct {
	connections map[int32][]*websocket.Conn
	mu          sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		connections: make(map[int32][]*websocket.Conn),
	}
}

func (s *WebSocketService) AddConnection(userID int32, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[userID] = append(s.connections[userID], conn)
}

func (s *WebSocketService) RemoveConnection(userID int32, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conns := s.connections[userID]
	for i, c := range conns {
		if c == conn {
			s.connections[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

func (s *WebSocketService) SendNotification(userID int32, notification WSNotification) {
	s.mu.RLock()
	conns := s.connections[userID]
	s.mu.RUnlock()

	msg := WSMessage{
		Type:    "notification",
		Payload: notification,
	}

	for _, conn := range conns {
		if err := conn.WriteJSON(msg); err != nil {
			s.RemoveConnection(userID, conn)
		}
	}

	// Отправляем событие для обновления бейджа
	triggerMsg := WSMessage{
		Type:    "trigger",
		Payload: "newNotification",
	}
	for _, conn := range conns {
		if err := conn.WriteJSON(triggerMsg); err != nil {
			s.RemoveConnection(userID, conn)
		}
	}
}

func (s *WebSocketService) CloseUserConnections(userID int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, conn := range s.connections[userID] {
		conn.Close()
	}
	s.connections[userID] = nil
}
