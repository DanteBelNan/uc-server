package domain

// Client representa a un usuario o dispositivo conectado a una sala.
// Contiene la informacion basica necesaria para identificar y comunicarse con el extremo.
type Client struct {
    ID   string
    Room string
}

// NewClient crea una nueva instancia de un cliente con un ID unico.
// id: identificador unico del cliente.
// room: identificador de la sala a la que pertenece.
func NewClient(id string, room string) *Client {
    return &Client{
        ID:   id,
        Room: room,
    }
}
