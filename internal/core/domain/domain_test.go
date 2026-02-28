package domain

import "testing"

// TestRoom_AddRemoveClient valida que la sala gestione correctamente sus clientes.
func TestRoom_AddRemoveClient(t *testing.T) {
    room := NewRoom("test-room")
    client := NewClient("client-1", "test-room")

    // Validacion inicial
    if count := room.GetClientsCount(); count != 0 {
        t.Errorf("Se esperaba 0 clientes, se obtuvo %d", count)
    }

    // Agregar cliente
    room.AddClient(client)
    if count := room.GetClientsCount(); count != 1 {
        t.Errorf("Se esperaba 1 cliente, se obtuvo %d", count)
    }

    // Remover cliente
    room.RemoveClient("client-1")
    if count := room.GetClientsCount(); count != 0 {
        t.Errorf("Se esperaba 0 clientes tras eliminar, se obtuvo %d", count)
    }
}

// TestMessageCreation valida que los mensajes se creen correctamente con sus atributos.
func TestMessageCreation(t *testing.T) {
    msgType := TypeUpdate
    payload := "contenido de prueba"
    sender := "client-99"

    msg := NewMessage(msgType, payload, sender)

    if msg.Type != msgType {
        t.Errorf("Tipo incorrecto: se esperaba %s, se obtuvo %s", msgType, msg.Type)
    }
    if msg.Payload != payload {
        t.Errorf("Payload incorrecto: se esperaba %s, se obtuvo %s", payload, msg.Payload)
    }
    if msg.Sender != sender {
        t.Errorf("Sender incorrecto: se esperaba %s, se obtuvo %s", sender, msg.Sender)
    }
}
