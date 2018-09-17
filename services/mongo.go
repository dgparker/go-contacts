package services

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoService implements ContactService
type MongoService struct {
	DBName     string
	Collection string
	Session    *mgo.Session
}

// NewMongoService returns an instance of MongoService
func NewMongoService(uri string, name string, collection string) *MongoService {
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}

	return &MongoService{
		DBName:     name,
		Collection: collection,
		Session:    session,
	}
}

// AllEntries returns all entries in the database
func (s *MongoService) AllEntries() ([]*Entry, error) {
	entries := []*Entry{}
	err := s.Session.DB(s.DBName).C(s.Collection).Find(bson.M{}).All(&entries)
	return entries, err
}

// EntryByID returns an entry for the given ID
func (s *MongoService) EntryByID(id string) (*Entry, error) {
	valid := bson.IsObjectIdHex(id)
	if !valid {
		return nil, ErrInvalidID
	}

	entry := &Entry{}
	err := s.Session.DB(s.DBName).C(s.Collection).Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&entry)
	return entry, err
}

// AddEntry adds a new entry into the database
func (s *MongoService) AddEntry(entry *Entry) (*Entry, error) {
	entry.Email = strings.ToLower(entry.Email)
	result := s.checkIfEmailExists(entry)
	if result != nil {
		return nil, result
	}

	err := s.Session.DB(s.DBName).C(s.Collection).Insert(entry)
	if err != nil {
		return nil, err
	}

	newEntry := &Entry{}
	err = s.Session.DB(s.DBName).C(s.Collection).Find(bson.M{"email": entry.Email}).One(&newEntry)
	return newEntry, err
}

// UpdateEntry updates an entry in the database
func (s *MongoService) UpdateEntry(entry *Entry) error {
	entry.Email = strings.ToLower(entry.Email)
	result := s.checkIfEmailExists(entry)
	if result != nil {
		return result
	}

	selector := bson.M{
		"_id": entry.ID,
	}
	err := s.Session.DB(s.DBName).C(s.Collection).Update(selector, entry)
	return err
}

// DeleteEntryByID removes an Entry from the Session.DB(s.DBName) using _id as the selector
func (s *MongoService) DeleteEntryByID(id string) error {
	valid := bson.IsObjectIdHex(id)
	if !valid {
		return ErrInvalidID
	}

	err := s.Session.DB(s.DBName).C(s.Collection).Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}

// EntriesToCSV retrieves all entries in the addressbook and writes them to a CSV file
func (s *MongoService) EntriesToCSV() (*os.File, error) {
	entries := []*Entry{}
	err := s.Session.DB(s.DBName).C(s.Collection).Find(bson.M{}).All(&entries)
	if err != nil {
		return nil, err
	}

	// write to temp file
	entriesFile, err := ioutil.TempFile(os.TempDir(), "tmp.*.csv")
	if err != nil {
		return nil, err
	}
	err = gocsv.MarshalFile(&entries, entriesFile)
	if err != nil {
		return nil, err
	}

	return entriesFile, nil
}

// CSVToEntries accepts entries in CSV format and attempts to insert them into the Session.DB(s.DBName)
func (s *MongoService) CSVToEntries(file *os.File) ([]*Entry, error) {
	entries := []*Entry{}
	failedEntries := []*Entry{}
	entryFile, err := os.Open(file.Name())
	if err != nil {
		return nil, err
	}
	defer func() {
		entryFile.Close()
		os.Remove(file.Name())
	}()

	err = gocsv.UnmarshalFile(entryFile, &entries)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		// If the ID field is provided update instead of add
		// error handling and logic could be further refined here depending on business requirements
		if entry.ID != "" {
			err = s.UpdateEntry(entry)
			if err != nil {
				switch err {
				case ErrEmailExists:
					failedEntries = append(failedEntries, entry)
				default:
					return failedEntries, err
				}

			}
		} else {
			_, err = s.AddEntry(entry)
			if err != nil {
				switch err {
				case ErrEmailExists:
					failedEntries = append(failedEntries, entry)
				default:
					return failedEntries, err
				}
			}
		}
	}

	if len(failedEntries) > 0 {
		err = ErrEmailExists
	}

	return failedEntries, err
}

// ErrEmailExists message for entry exists in Session.DB(s.DBName) on AddEntry
var ErrEmailExists = errors.New("Email already exists")
var errNotFound = errors.New("not found")

func (s *MongoService) checkIfEmailExists(entry *Entry) error {
	selector := bson.M{
		"email": entry.Email,
	}

	result := &Entry{}
	err := s.Session.DB(s.DBName).C(s.Collection).Find(selector).One(&result)
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
			return nil
		default:
			return err
		}
	}
	// use our structs validate function here to easily check to see
	// if our db returned a document if it's valid then return email exists error
	err = result.Validate()
	if err == nil && result.ID != entry.ID {
		return ErrEmailExists
	}

	return nil
}
