package services

import (
    "errors"
    "sync"

    "github.com/DanteBelNan/uc-server/internal/core/domain"
)

// RoomService gestiona el ciclo de vida de las salas y la comunicacion entre clientes.
// Utiliza un mapa en memoria protegido por un Mutex para garantizar la seguridad concurrente.
type RoomService struct {
    rooms map[string]*domain.Room
    mu    sync.RWMutex
}

// NewRoomService crea una nueva instancia del servicio de salas.
func NewRoomService() *RoomService {
    return &RoomService{
        rooms: make(map[string]*domain.Room),
    }
}

// GetOrCreateRoom busca una sala por su ID o la crea si no existe.
// roomID: identificador unico de la sala.
func (s *RoomService) GetOrCreateRoom(roomID string) *domain.Room {
    s.mu.Lock()
    defer s.mu.Unlock()

    if room, exists := s.rooms[roomID]; exists {
        return room
    }

    newRoom := domain.NewRoom(roomID)
    s.rooms[roomID] = newRoom
    return newRoom
}

// JoinRoom agrega un cliente a una sala especifica.
// roomID: ID de la sala a la que unirse.
// client: puntero al cliente que desea unirse.
func (s *RoomService) JoinRoom(roomID string, client *domain.Client) error {
    if roomID == "" {
        return errors.New("el ID de la sala no puede estar vacio")
    }

    room := s.GetOrCreateRoom(roomID)
    room.AddClient(client)
    return nil
}

// LeaveRoom elimina a un cliente de una sala y limpia la sala si queda vacia.
// roomID: ID de la sala.
// clientID: ID del cliente que abandona.
func (s *RoomService) LeaveRoom(roomID string, clientID string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    room, exists := s.rooms[roomID]
    if !exists {
        return
    }

    room.RemoveClient(clientID)

    // Si la sala no tiene mas clientes, se elimina para liberar memoria.
    if room.GetClientsCount() == 0 {
        delete(s.rooms, roomID)
    }
}

// BroadcastMessage envia un mensaje a todos los clientes de una sala excepto al emisor.
// roomID: ID de la sala donde se enviara el mensaje.
// msg: el mensaje a retransmitir.
// El ruteo fisico de los mensajes se delegara a los adaptadores de salida en pasos posteriores.
func (s *RoomService) BroadcastMessage(roomID string, msg domain.Message) ([]*domain.Client, error) {
    s.mu.RLock()
    room, exists := s.rooms[roomID]
    s.mu.RUnlock()

    if !exists {
        return nil, errors.New("la sala no existe")
    }

    var recipients []*domain.Client
    room.Mu.RLock()
    defer room.Mu.RUnlock()

    for _, client := range room.Clients {
        if client.ID != msg.Sender {
            recipients = append(recipients, client)
        }
    }

    return recipients, nil
}
