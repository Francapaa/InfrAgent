package ws

import (
	"bytes"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("[WebSocket] CheckOrigin - Origin: %s", r.Header.Get("Origin"))
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	// Recovery de panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WebSocket] PANIC en readPump: %v\n%s", r, debug.Stack())
			c.conn.Close()
		}

		log.Println("[WebSocket] Cerrando readPump y desregistrando cliente")

		// Desregistrar del hub (no bloqueante)
		select {
		case c.hub.unregister <- c:
			log.Println("[WebSocket] Cliente desregistrado exitosamente")
		case <-time.After(5 * time.Second):
			log.Println("[WebSocket] Timeout al desregistrar cliente")
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		log.Println("[WebSocket] Pong recibido")
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	log.Println("[WebSocket] readPump iniciado")

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] Error inesperado en conexión: %v", err)
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("[WebSocket] Cliente cerró conexión normalmente (1000)")
			} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Println("[WebSocket] Cliente se fue (1001)")
			} else {
				log.Printf("[WebSocket] Conexión cerrada - Error: %v (Tipo: %T)", err, err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
		log.Printf("[WebSocket] Mensaje recibido: %s", string(message))

		// Enviar al hub de forma no bloqueante
		select {
		case c.hub.broadcast <- message:
		default:
			log.Println("[WebSocket] Canal broadcast lleno, mensaje descartado")
		}
	}
}

func (c *Client) writePump() {
	// Recovery de panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WebSocket] PANIC en writePump: %v\n%s", r, debug.Stack())
		}
		log.Println("[WebSocket] Cerrando writePump")
	}()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	log.Println("[WebSocket] writePump iniciado")

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// El hub cerró el canal
				log.Println("[WebSocket] Canal de envío cerrado por el hub")
				// Enviar Close Frame apropiado
				c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("[WebSocket] Error al obtener writer: %v", err)
				return
			}

			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newLine)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Printf("[WebSocket] Error al cerrar writer: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[WebSocket] Error al enviar ping: %v", err)
				return
			}
			log.Println("[WebSocket] Ping enviado")
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf("[WebSocket] Nueva solicitud de conexión desde: %s", r.RemoteAddr)
	log.Printf("[WebSocket] Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s", name, value)
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WebSocket] FALLO EL UPGRADER: %v", err)
		return
	}

	log.Printf("[WebSocket] WebSocket upgrade exitoso para %s", r.RemoteAddr)

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	// Registrar cliente de forma segura
	select {
	case hub.register <- client:
		log.Printf("[WebSocket] Cliente enviado a registro")
	case <-time.After(5 * time.Second):
		log.Println("[WebSocket] Timeout al registrar cliente")
		conn.Close()
		return
	}

	go client.writePump()
	go client.readPump()

	log.Printf("[WebSocket] Goroutines iniciadas para cliente %s", r.RemoteAddr)
}
