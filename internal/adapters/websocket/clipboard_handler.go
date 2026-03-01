package websocket

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/DanteBelNan/uc-server/internal/core/domain"
    "github.com/DanteBelNan/uc-server/internal/core/services"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        // En una version posterior se deberia restringir el origen por seguridad.
        return true
    },
}

// ClipboardHandler maneja las conexiones WebSocket y las integra con el RoomService.
// Mantiene un mapa de conexiones activas para poder enviar mensajes a clientes especificos.
type ClipboardHandler struct {
    roomService *services.RoomService
    conns       map[string]*websocket.Conn
    mu          sync.RWMutex
}

// NewClipboardHandler crea una nueva instancia del manejador de sockets.
// roomService: el servicio de negocio que gestiona las salas.
func NewClipboardHandler(roomService *services.RoomService) *ClipboardHandler {
    return &ClipboardHandler{
        roomService: roomService,
        conns:       make(map[string]*websocket.Conn),
    }
}

// HandleConnection gestiona el upgrade de HTTP a WebSocket y el ciclo de vida de la conexion.
func (h *ClipboardHandler) HandleConnection(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Error al hacer upgrade a WebSocket: %v", err)
        return
    }
    defer conn.Close()

    var currentClientID string
    var currentRoomID string

    // Bucle principal para leer mensajes del cliente.
    for {
        var msg domain.Message
        if err := conn.ReadJSON(&msg); err != nil {
            log.Printf("Conexion cerrada o error de lectura: %v", err)
            if currentClientID != "" && currentRoomID != "" {
                h.handleDisconnect(currentRoomID, currentClientID)
            }
            break
        }

        switch msg.Type {
        case domain.TypeJoin:
            currentClientID = msg.Sender
            currentRoomID = msg.Payload
            h.handleJoin(currentRoomID, currentClientID, conn)

        case domain.TypeUpdate:
            if currentRoomID != "" {
                h.handleUpdate(currentRoomID, msg)
            }
        }
    }
}

// handleJoin procesa la solicitud de unirse a una sala.
func (h *ClipboardHandler) handleJoin(roomID string, clientID string, conn *websocket.Conn) {
    client := domain.NewClient(clientID, roomID)
    if err := h.roomService.JoinRoom(roomID, client); err != nil {
        log.Printf("Error al unir cliente a sala: %v", err)
        return
    }

    h.mu.Lock()
    h.conns[clientID] = conn
    h.mu.Unlock()

    log.Printf("Cliente %s unido a la sala %s", clientID, roomID)
    h.broadcastUserList(roomID)
}

// handleUpdate retransmite las actualizaciones del portapapeles a otros clientes.
func (h *ClipboardHandler) handleUpdate(roomID string, msg domain.Message) {
    recipients, err := h.roomService.BroadcastMessage(roomID, msg)
    if err != nil {
        log.Printf("Error al obtener destinatarios para broadcast: %v", err)
        return
    }

    h.mu.RLock()
    defer h.mu.RUnlock()

    for _, recipient := range recipients {
        if conn, exists := h.conns[recipient.ID]; exists {
            if err := conn.WriteJSON(msg); err != nil {
                log.Printf("Error enviando mensaje a %s: %v", recipient.ID, err)
            }
        }
    }
}

// handleDisconnect limpia los recursos del cliente al desconectarse.
func (h *ClipboardHandler) handleDisconnect(roomID string, clientID string) {
    h.roomService.LeaveRoom(roomID, clientID)
    h.mu.Lock()
    delete(h.conns, clientID)
    h.mu.Unlock()
    log.Printf("Cliente %s desconectado de la sala %s", clientID, roomID)
    h.broadcastUserList(roomID)
}

// broadcastUserList envia a todos los clientes de una sala la lista actual de miembros.
func (h *ClipboardHandler) broadcastUserList(roomID string) {
    userList, err := h.roomService.GetRoomUsers(roomID)
    if err != nil {
        return
    }

    data, _ := json.Marshal(userList)
    msg := domain.NewMessage(domain.TypeUserList, string(data), "SERVER")

    h.mu.RLock()
    defer h.mu.RUnlock()

    for _, userID := range userList {
        if conn, exists := h.conns[userID]; exists {
            conn.WriteJSON(msg)
        }
    }
}
