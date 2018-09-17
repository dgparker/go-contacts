package services

import (
	"errors"
	"os"

	"gopkg.in/mgo.v2/bson"
)

var (
	// ErrInvalidName message for invalid name field
	ErrInvalidName = errors.New("Name field cannot be empty or null")

	// ErrInvalidEmail message for invalid email field
	ErrInvalidEmail = errors.New("Email field cannot be empty or null")

	// ErrInvalidPhone message for invalid phone field
	ErrInvalidPhone = errors.New("Phone field cannot be empty or null")

	// ErrInvalidID message for invalid ID field
	ErrInvalidID = errors.New("ID value is invalid")
)

// ContactService defines the methods necessary to implement a contact service
type ContactService interface {
	AllEntries() ([]*Entry, error)
	EntryByID(id string) (*Entry, error)
	AddEntry(entry *Entry) (*Entry, error)
	UpdateEntry(entry *Entry) error
	DeleteEntryByID(id string) error
	EntriesToCSV() (*os.File, error)
	CSVToEntries(*os.File) ([]*Entry, error)
}

// Entry defines what an addressbook entry should contain
type Entry struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	Email     string        `json:"email" bson:"email"`
	Phone     string        `json:"phone" bson:"phone"`
}

// Validation interface is implemented by the Entry struct and used to validate "struct" specific fields.
type Validation interface {
	Validate() error
}

// Validate ensures data set doesn't contain empty or null values
//
// this could further be refined to determine whether or not each field
// is in a valid format i.e. email and phone numbers or length of names are < a specified amount
func (e Entry) Validate() error {
	if e.FirstName == "" {
		return ErrInvalidName
	}
	if e.LastName == "" {
		return ErrInvalidName
	}
	if e.Email == "" {
		return ErrInvalidEmail
	}
	if e.Phone == "" {
		return ErrInvalidPhone
	}
	return nil
}
