package domain

import "sync"

// Room representa una sala de sincronizacion donde multiples clientes estan suscritos.
// Gestiona la lista de clientes activos de forma segura para acceso concurrente.
type Room struct {
    ID      string
    Clients map[string]*Client
    Mu      sync.RWMutex
}

// NewRoom crea una nueva instancia de sala inicializada.
// id: identificador unico de la sala.
func NewRoom(id string) *Room {
    return &Room{
        ID:      id,
        Clients: make(map[string]*Client),
    }
}

// AddClient agrega un cliente a la sala de forma segura.
// client: puntero a la instancia del cliente a agregar.
func (r *Room) AddClient(client *Client) {
    r.Mu.Lock()
    defer r.Mu.Unlock()
    r.Clients[client.ID] = client
}

// RemoveClient elimina un cliente de la sala por su ID.
// clientID: identificador unico del cliente a remover.
func (r *Room) RemoveClient(clientID string) {
    r.Mu.Lock()
    defer r.Mu.Unlock()
    delete(r.Clients, clientID)
}

// GetClientsCount devuelve la cantidad de clientes conectados actualmente.
func (r *Room) GetClientsCount() int {
    r.Mu.RLock()
    defer r.Mu.RUnlock()
    return len(r.Clients)
}
