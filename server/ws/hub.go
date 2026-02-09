package ws

import (
	"log"
	"time"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("[Hub] Cliente registrado. Total: %d clientes", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				log.Printf("[Hub] Cliente desregistrado. Total: %d clientes", len(h.clients))

				// Cerrar el canal de envío para que writePump se entere
				// Esto causará que writePump reciba !ok y envíe Close Frame
				close(client.send)
			}

		case message := <-h.broadcast:
			log.Printf("[Hub] Broadcast de mensaje a %d clientes", len(h.clients))
			for client := range h.clients {
				select {
				case client.send <- message:
					// Éxito
				default:
					// Canal lleno o bloqueado, remover cliente
					log.Println("[Hub] Canal de cliente lleno, desregistrando")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

var hub *Hub

func init() {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256), // Buffer para evitar bloqueos
		register:   make(chan *Client, 10), // Buffer para registro
		unregister: make(chan *Client, 10), // Buffer para desregistro
	}
	go hub.run()
	log.Println("[Hub] Inicializado y corriendo")
}

func GetHub() *Hub {
	return hub
}

// BroadcastMessage envía un mensaje a todos los clientes conectados
func BroadcastMessage(message []byte) {
	select {
	case hub.broadcast <- message:
		// Éxito
	case <-time.After(5 * time.Second):
		log.Println("[Hub] Timeout al enviar mensaje al broadcast")
	}
}
