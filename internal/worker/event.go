package worker

type UserCreatedEvent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
