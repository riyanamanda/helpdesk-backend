package mailer

import (
	"fmt"
)

func NewTicketMessage(ticketID int64, title string, description string, submitterName string, adminEmails []string) (Message, error) {
	data := TicketCreatedTemplateData{
		TicketID:    ticketID,
		Title:       title,
		Description: description,
		SubmittedBy: submitterName,
	}

	htmlBody, err := renderTicketCreatedHTML(data)
	if err != nil {
		return Message{}, fmt.Errorf("render ticket created HTML: %w", err)
	}

	return Message{
		To:      adminEmails[0],
		CC:      adminEmails[1:],
		Subject: fmt.Sprintf("New Ticket #%d: %s", ticketID, title),
		Body:    htmlBody,
	}, nil
}
