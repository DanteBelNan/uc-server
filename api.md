# Documentacion de la API de Universal Clipboard (v0.1.x)

Esta documentacion detalla el protocolo de comunicacion entre los clientes y el servidor central (`uc-server`) para la sincronizacion del portapapeles.

## 📡 Endpoint WebSocket
El servidor expone un unico punto de entrada para WebSockets:
- **URL:** `ws://<host>:<port>/ws`
- **Protocolo:** JSON sobre WebSockets.

---

## 🏗 Flujo de Comunicacion

### 1. Conexion y Union a Sala (JOIN)
Todo cliente debe enviar un mensaje de tipo `JOIN` inmediatamente despues de establecer la conexion para empezar a recibir actualizaciones.

**Mensaje de Cliente:**
```json
{
    "type": "JOIN",
    "payload": "mi-sala-segura",
    "sender_id": "dispositivo-uuid-123"
}
```
- `payload`: El nombre unico de la sala.
- `sender_id`: Identificador unico del dispositivo.

### 2. Sincronizacion del Portapapeles (UPDATE)
Cuando un cliente detecta un cambio en el portapapeles local, debe enviarlo a la sala. El servidor lo retransmitira a todos los demas suscriptores activos.

**Mensaje de Cliente:**
```json
{
    "type": "UPDATE",
    "payload": "Texto copiado en el portapapeles",
    "sender_id": "dispositivo-uuid-123"
}
```

**Recepcion en Otros Clientes:**
Los otros dispositivos de la sala recibiran el mismo JSON. El cliente receptor debe verificar el `sender_id` para evitar procesar sus propios mensajes (aunque el servidor ya filtra el eco).

---

## 📝 Tipos de Mensajes Soportados

| Tipo | Origen | Descripcion |
| :--- | :--- | :--- |
| `JOIN` | Cliente | Solicita unirse a una sala especifica. |
| `UPDATE` | Cliente/Servidor | Envia o recibe contenido del portapapeles. |
| `ERROR` | Servidor | Notifica fallos en el protocolo o en la sala. |
| `HEARTBEAT`| Ambos | (Opcional) Mantiene la conexion activa ante inactividad. |

---

## ⚠️ Consideraciones de la Version 0.1.x
- El servidor actualmente no requiere contraseña para unirse a una sala (se implementara en v0.2.x/v0.3.x).
- Los datos viajan en **texto plano** a traves del socket. No enviar informacion sensible hasta la implementacion de E2EE en la Fase 3.
