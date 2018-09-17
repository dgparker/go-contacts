package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"

	"github.com/dgparker/go-contacts/services"
)

type testDataService struct{}

func (tds *testDataService) AllEntries() ([]*services.Entry, error) {
	return []*services.Entry{}, nil
}

func (tds *testDataService) EntryByID(id string) (*services.Entry, error) {
	return &services.Entry{}, nil
}

func (tds *testDataService) AddEntry(entry *services.Entry) (*services.Entry, error) {
	return &services.Entry{}, nil
}

func (tds *testDataService) UpdateEntry(entry *services.Entry) error {
	return nil
}

func (tds *testDataService) DeleteEntryByID(id string) error {
	return nil
}

func (tds *testDataService) EntriesToCSV() (*os.File, error) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "tmp.*.csv")
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

func (tds *testDataService) CSVToEntries(file *os.File) ([]*services.Entry, error) {
	return []*services.Entry{}, nil
}

func TestHandleGetEntries(t *testing.T) {
	testContactHandler := NewContactHandler(&testDataService{})

	r := httprouter.New()
	r.GET("/entries", testContactHandler.HandleGetEntries)

	req, err := http.NewRequest("GET", "/entries", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code: %d to equal: %d", status, http.StatusOK)
	}
}

func TestHandleGetEntryByID(t *testing.T) {
	testContactHandler := NewContactHandler(&testDataService{})

	r := httprouter.New()
	r.GET("/entries/:id", testContactHandler.HandleGetEntryByID)

	req, err := http.NewRequest("GET", "/entries/1234", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code: %d to equal: %d", status, http.StatusOK)
	}
}

var testEntry = &services.Entry{
	ID:        bson.NewObjectId(),
	FirstName: "Tester",
	LastName:  "McTesterson",
	Email:     "tester.mctesterson@codecave.app",
	Phone:     "8675309",
}

func TestHandlePostEntry(t *testing.T) {
	testJSON, err := json.Marshal(testEntry)
	if err != nil {
		t.Error(err)
	}

	testContactHandler := NewContactHandler(&testDataService{})

	r := httprouter.New()
	r.POST("/entry", testContactHandler.HandlePostEntry)

	req, err := http.NewRequest("POST", "/entry", bytes.NewBuffer(testJSON))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code: %d to equal: %d", status, http.StatusOK)
	}
}

func TestHandlePutEntry(t *testing.T) {
	testJSON, err := json.Marshal(testEntry)
	if err != nil {
		t.Error(err)
	}

	testContactHandler := NewContactHandler(&testDataService{})

	r := httprouter.New()
	r.PUT("/entry", testContactHandler.HandlePutEntry)

	req, err := http.NewRequest("PUT", "/entry", bytes.NewBuffer(testJSON))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Expected status code: %d to equal: %d", status, http.StatusAccepted)
	}
}

func TestHandleDeleteEntryByID(t *testing.T) {
	testContactHandler := NewContactHandler(&testDataService{})

	r := httprouter.New()
	r.DELETE("/entry/:id", testContactHandler.HandleDeleteEntryByID)

	req, err := http.NewRequest("DELETE", "/entry/1234", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code: %d to equal: %d", status, http.StatusOK)
	}
}

func TestHandleGetCSV(t *testing.T) {
	testContactHandler := NewContactHandler(&testDataService{})

	r := httprouter.New()
	r.GET("/csv/entry", testContactHandler.HandleGetCSV)

	req, err := http.NewRequest("GET", "/csv/entry", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code: %d to equal: %d", status, http.StatusOK)
	}
}

func createCSVFormFile(w *multipart.Writer, filename string) (io.Writer, error) {
	mh := make(textproto.MIMEHeader)
	mh.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	mh.Set("Content-Type", "text/csv")
	return w.CreatePart(mh)
}

func TestHandlePostCSV(t *testing.T) {
	testContactHandler := NewContactHandler(&testDataService{})

	path := filepath.Join("../test/csvtest.csv")
	testCSV, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}
	defer testCSV.Close()

	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)
	formFile, _ := createCSVFormFile(bodyWriter, "csvtest.csv")
	io.Copy(formFile, testCSV)
	bodyWriter.Close()

	r := httprouter.New()
	r.POST("/csv/entry", testContactHandler.HandlePostCSV)

	req, err := http.NewRequest("POST", "/csv/entry", bodyBuffer)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Expected status code: %d to equal: %d", status, http.StatusAccepted)
	}
}
