package domain

// MessageType define los tipos de mensajes permitidos en el protocolo de la aplicacion.
type MessageType string

const (
    // TypeJoin se utiliza cuando un cliente solicita unirse a una sala.
    TypeJoin      MessageType = "JOIN"
    // TypeUpdate se utiliza para enviar actualizaciones del portapapeles.
    TypeUpdate    MessageType = "UPDATE"
    // TypeError representa un mensaje de error enviado desde el servidor.
    TypeError     MessageType = "ERROR"
    // TypeUserList se utiliza para enviar la lista de usuarios conectados en la sala.
    TypeUserList  MessageType = "USER_LIST"
    // TypeHeartbeat se utiliza para mantener activa la conexion.
    TypeHeartbeat MessageType = "HEARTBEAT"
)

// Message es la estructura estandar para la comunicacion entre clientes y servidor.
// Se serializa a JSON para su transmision a traves de WebSockets.
type Message struct {
    Type    MessageType `json:"type"`
    Payload string      `json:"payload"`
    Sender  string      `json:"sender_id,omitempty"`
}

// NewMessage crea una nueva instancia de mensaje para ser enviada.
// msgType: el tipo de mensaje segun el protocolo.
// payload: el contenido del mensaje (texto plano o datos serializados).
// sender: ID del cliente que envia el mensaje.
func NewMessage(msgType MessageType, payload string, sender string) Message {
    return Message{
        Type:    msgType,
        Payload: payload,
        Sender:  sender,
    }
}
