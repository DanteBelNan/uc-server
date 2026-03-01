package services

import (
    "testing"

    "github.com/DanteBelNan/uc-server/internal/core/domain"
)

// TestRoomService_JoinLeaveRoom valida que los clientes puedan unirse y abandonar salas correctamente.
func TestRoomService_JoinLeaveRoom(t *testing.T) {
    service := NewRoomService()
    roomID := "sala-test"
    client := domain.NewClient("client-1", roomID)

    // Unir cliente
    err := service.JoinRoom(roomID, client)
    if err != nil {
        t.Errorf("Error al unirse a la sala: %v", err)
    }

    room := service.GetOrCreateRoom(roomID)
    if room.GetClientsCount() != 1 {
        t.Errorf("Se esperaba 1 cliente en la sala, se obtuvo %d", room.GetClientsCount())
    }

    // Abandonar cliente
    service.LeaveRoom(roomID, "client-1")
    
    // Validar limpieza de sala vacia
    service.mu.RLock()
    _, exists := service.rooms[roomID]
    service.mu.RUnlock()
    if exists {
        t.Errorf("La sala deberia haber sido eliminada tras quedar vacia")
    }
}

// TestRoomService_BroadcastMessage valida que el mensaje se distribuya a los destinatarios correctos.
func TestRoomService_BroadcastMessage(t *testing.T) {
    service := NewRoomService()
    roomID := "sala-broadcast"
    
    c1 := domain.NewClient("c1", roomID)
    c2 := domain.NewClient("c2", roomID)
    c3 := domain.NewClient("c3", roomID)

    service.JoinRoom(roomID, c1)
    service.JoinRoom(roomID, c2)
    service.JoinRoom(roomID, c3)

    msg := domain.NewMessage(domain.TypeUpdate, "hola mundo", "c1")
    recipients, err := service.BroadcastMessage(roomID, msg)

    if err != nil {
        t.Errorf("Error al intentar broadcast: %v", err)
    }

    // Deberia haber 2 destinatarios (c2 y c3), excluyendo al emisor (c1)
    if len(recipients) != 2 {
        t.Errorf("Se esperaban 2 destinatarios, se obtuvieron %d", len(recipients))
    }

    // Verificar exclusion del emisor
    for _, r := range recipients {
        if r.ID == "c1" {
            t.Errorf("El emisor (c1) no deberia estar en la lista de destinatarios")
        }
    }
}
