// Package models consists of data models.
package models

import (
	"time"
)

// User represents the application user.
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	PassHash  []byte    `json:"-" db:"hashed_password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Password represents the login-password pair.
type Password struct {
	ID        int       `json:"id" db:"id"`
	OwnerID   int       `json:"owner_id" db:"owner_id"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"password" db:"password"`
	Metadata  string    `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ToEntity transforms password model to dto entity.
func (p *Password) ToEntity() *Entity {
	return &Entity{
		ID:       p.ID,
		OwnerID:  p.OwnerID,
		Type:     TypePassword,
		Login:    p.Login,
		Password: p.Password,
		Metadata: p.Metadata,
	}
}

// Card represents the bank card.
type Card struct {
	ID        int       `json:"id" db:"id"`
	OwnerID   int       `json:"owner_id" db:"owner_id"`
	Number    string    `json:"number" db:"number"`
	CVC       string    `json:"cvc" db:"cvc"`
	Owner     string    `json:"owner" db:"owner"`
	Date      string    `json:"date" db:"date"`
	Metadata  string    `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ToEntity transforms card model to dto entity.
func (c *Card) ToEntity() *Entity {
	return &Entity{
		ID:         c.ID,
		OwnerID:    c.OwnerID,
		Type:       TypeCard,
		CardNumber: c.Number,
		CardOwner:  c.Owner,
		CardCVC:    c.CVC,
		CardExp:    c.Date,
		Metadata:   c.Metadata,
	}
}

// Text represents the text information.
type Text struct {
	ID        int       `json:"id" db:"id"`
	OwnerID   int       `json:"owner_id" db:"owner_id"`
	Text      string    `json:"text" db:"text"`
	Metadata  string    `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ToEntity transforms text model to dto entity.
func (t *Text) ToEntity() *Entity {
	return &Entity{
		ID:      t.ID,
		OwnerID: t.OwnerID,
		Type:    TypeText,
		Text:    t.Text,
	}
}

// File represents the text uploaded file.
type File struct {
	ID        int       `json:"id" db:"id"`
	OwnerID   int       `json:"owner_id" db:"owner_id"`
	FileName  string    `json:"filename" db:"filename"`
	Metadata  string    `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ToEntity transforms file model to dto entity.
func (f *File) ToEntity() *Entity {
	return &Entity{
		ID:       f.ID,
		OwnerID:  f.OwnerID,
		Type:     TypeFile,
		Metadata: f.Metadata,
		Filename: f.FileName,
	}
}

// EntityType is the Entity type.
type EntityType string

const (
	TypePassword = EntityType("PASSWORD")
	TypeCard     = EntityType("CARD")
	TypeText     = EntityType("TEXT")
	TypeFile     = EntityType("FILE")
)

// Entity represents the entity DTO.
type Entity struct {
	ID         int
	OwnerID    int
	Type       EntityType
	Login      string
	Password   string
	CardNumber string
	CardOwner  string
	CardCVC    string
	CardExp    string
	Text       string
	Filename   string
	Metadata   string
}

// ToPassword transforms entity to the password model.
func (e *Entity) ToPassword() *Password {
	return &Password{
		ID:       e.ID,
		OwnerID:  e.OwnerID,
		Login:    e.Login,
		Password: e.Password,
		Metadata: e.Metadata,
	}
}

// ToCard transforms entity to the card model.
func (e *Entity) ToCard() *Card {
	return &Card{
		ID:       e.ID,
		OwnerID:  e.OwnerID,
		Number:   e.CardNumber,
		CVC:      e.CardCVC,
		Owner:    e.CardOwner,
		Date:     e.CardExp,
		Metadata: e.Metadata,
	}
}

// ToText transforms entity to the text model.
func (e *Entity) ToText() *Text {
	return &Text{
		ID:       e.ID,
		OwnerID:  e.OwnerID,
		Text:     e.Text,
		Metadata: e.Metadata,
	}
}
