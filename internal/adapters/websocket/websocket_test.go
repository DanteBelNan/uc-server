package websocket

import (
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/DanteBelNan/uc-server/internal/core/domain"
    "github.com/DanteBelNan/uc-server/internal/core/services"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

// TestWebSocket_Integration verifica el flujo completo de sincronizacion entre dos clientes.
func TestWebSocket_Integration(t *testing.T) {
    // Configurar el entorno de prueba
    gin.SetMode(gin.TestMode)
    roomService := services.NewRoomService()
    handler := NewClipboardHandler(roomService)

    router := gin.New()
    router.GET("/ws", func(c *gin.Context) {
        handler.HandleConnection(c)
    })

    server := httptest.NewServer(router)
    defer server.Close()

    // Crear la URL de WebSocket
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

    // Cliente 1: Se conecta y se une a la sala "sala-1"
    c1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("Fallo al conectar cliente 1: %v", err)
    }
    defer c1.Close()

    joinMsg1 := domain.NewMessage(domain.TypeJoin, "sala-1", "client-1")
    if err := c1.WriteJSON(joinMsg1); err != nil {
        t.Fatalf("Error al enviar JOIN desde cliente 1: %v", err)
    }

    // Cliente 2: Se conecta y se une a la misma sala "sala-1"
    c2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("Fallo al conectar cliente 2: %v", err)
    }
    defer c2.Close()

    joinMsg2 := domain.NewMessage(domain.TypeJoin, "sala-1", "client-2")
    if err := c2.WriteJSON(joinMsg2); err != nil {
        t.Fatalf("Error al enviar JOIN desde cliente 2: %v", err)
    }

    // Pequena espera para asegurar el procesamiento de JOIN en el servidor
    time.Sleep(50 * time.Millisecond)

    // Cliente 1 envia una actualizacion del portapapeles
    payload := "contenido compartido"
    updateMsg := domain.NewMessage(domain.TypeUpdate, payload, "client-1")
    if err := c1.WriteJSON(updateMsg); err != nil {
        t.Fatalf("Error al enviar UPDATE desde cliente 1: %v", err)
    }

    // Cliente 2 deberia recibir la actualizacion
    var receivedMsg domain.Message
    c2.SetReadDeadline(time.Now().Add(1 * time.Second))
    if err := c2.ReadJSON(&receivedMsg); err != nil {
        t.Fatalf("Cliente 2 no recibio el mensaje: %v", err)
    }

    if receivedMsg.Payload != payload {
        t.Errorf("Contenido incorrecto: se esperaba %s, se obtuvo %s", payload, receivedMsg.Payload)
    }

    if receivedMsg.Sender != "client-1" {
        t.Errorf("Emisor incorrecto: se esperaba client-1, se obtuvo %s", receivedMsg.Sender)
    }
}
