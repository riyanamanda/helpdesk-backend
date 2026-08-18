package websocket

type Publisher interface {
	Publish(message Message)
}

type HubPublisher struct {
	hub *Hub
}

func NewHubPublisher(hub *Hub) Publisher {
	return &HubPublisher{
		hub: hub,
	}
}

func (p *HubPublisher) Publish(message Message) {
	p.hub.Broadcast(message)
}
