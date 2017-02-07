package tesis

//This file defines the Conversations
//carried by the user

type Conversation interface {
	Messages() []Message
	Participate(subject, recipient, body string)
}

type Message interface {
	Subject() string
	Recipient() string
	Remitent() string
	Body() string
}

type Credentials struct {
	user string
	pass string
}

type Portal interface {
	Auth(c Credentials) bool
	Conversate() []Conversation
}
