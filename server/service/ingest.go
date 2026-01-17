package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	models "server/model"
	"server/repositories"
	"server/utils"
)

type IngestHandler struct {
	events repositories.EventRepository
	agent  repositories.AgentStorage
	client repositories.ClientStorage
}

func NewIngestHandler(e repositories.EventRepository, a repositories.AgentStorage) *IngestHandler {
	return &IngestHandler{events: e, agent: a}
}

//RECIBIR LA API KEY PARA VER A QUE AGENTE PERTENECE
//PARSEAMOS EL EVENTO QUE NOS MANDA EL SDK
//AL EVENTO NUEVO LE PONEMOS AGENT ID Y CLIENT ID
//GUARDAMOS EL EVENTO EN LA BASE DE DATOS (POSTGRE)

func (IH *IngestHandler) NewEventInRequestService(ctx context.Context, apiKey string, event *models.Event) error {

	apiKeyHashed := utils.HashAPIKey(apiKey)

	agent, err := IH.agent.GetAgentByApiKey(ctx, apiKeyHashed)

	if err != nil {
		return errors.New("la api key no esta registrada")
	}

	event.AgentID = agent.ID
	event.ClientID = agent.ClientID

	if event.Severity == "critical" {
		go IH.sendUrgentNotification(agent.ClientID, event.Service) // disparar en segundo plano
	}

	return IH.events.CreateEvent(ctx, event)

}

func (IH *IngestHandler) sendUrgentNotification(clientId string, servicio string) {
	//LOGICA QUE VA A CONECTAR CON MAIL
	ctx := context.Background()
	client, err := IH.client.GetClient(ctx, clientId)
	if err != nil {
		log.Printf("[Error] No se pudo encontrar el cliente %s para notificar: %v", clientId, err)
		return
	}

	// 2. Configurar el mensaje
	subject := "⚠️ Acción de Emergencia Detectada - " + client.CompanyName
	body := fmt.Sprintf(
		"Hola %s,\n\nNuestro Agente detectó un fallo crítico en tu servicio: %s.\n"+
			"Se va a proceder a ejecutar una acción de recuperación automática (Restart/Scale).\n\n"+
			"Puedes monitorear esto en tiempo real desde tu Dashboard.",
		client.Nombre, servicio,
	)

	// 3. Enviar el Mail (Usando una función auxiliar)
	err = IH.sendMail(client.Email, subject, body)
	if err != nil {
		log.Printf("[Error] Falló el envío de mail a %s: %v", client.Email, err)
	} else {
		log.Printf("[Notificación] Mail enviado con éxito a %s", client.Email)
	}
}

func (IH *IngestHandler) sendMail(email, subject, body string) error {

	from := "franciscocaparruva@gmail.com"
	password := "Libertadores2018" // No es la contraseña real, es una 'App Password'
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Construir el mensaje con Headers correctos
	message := []byte(
		"I'M FRANCISCO CAPARRUVA CEO AND CTO OF INFRAGENT\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n")

	// Autenticación
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Envío de TLS
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, message)
	return err

}
