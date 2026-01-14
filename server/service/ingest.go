package service

import (
	"context"
	"errors"
	models "server/model"
	"server/repositories"
	"server/utils"
)

type IngestHandler struct {
	events repositories.EventRepository
	agent  repositories.AgentStorage
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
		go IH.sendUrgentNotification(agent.ClientID, event.Service)
	}

	return IH.events.CreateEvent(ctx, event)

}

func (IH *IngestHandler) sendUrgentNotification(clientId string, servicio string) {
	//LOGICA QUE VA A CONECTAR CON MAIL
}
