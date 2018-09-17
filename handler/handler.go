package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgparker/go-contacts/services"
	"github.com/julienschmidt/httprouter"
)

// decodeAndValidate is a simple helper function that decodes and validates struct fields
func decodeAndValidate(r *http.Request, v services.Validation) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	defer r.Body.Close()

	return v.Validate()
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}

// ContactHandler contains the methods for handling HTTP requests
// and exposes the ContactService
type ContactHandler struct {
	Service services.ContactService
	Logger  *log.Logger
}

// NewContactHandler returns a new ContactHandler
func NewContactHandler(svc services.ContactService) *ContactHandler {
	return &ContactHandler{
		Service: svc,
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
	}
}

// HandleGetEntries handles requests to GET /entry
func (h *ContactHandler) HandleGetEntries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.Service.AllEntries()
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	encodeJSON(w, res, h.Logger)
}

// HandleGetEntryByID handles requests to GET /entry/:id
func (h *ContactHandler) HandleGetEntryByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		Error(w, ErrNoIDParam, http.StatusBadRequest, h.Logger)
		return
	}

	res, err := h.Service.EntryByID(id)
	if err != nil {
		switch err {
		case services.ErrInvalidID:
			Error(w, err, http.StatusBadRequest, h.Logger)
			return
		default:
			Error(w, err, http.StatusOK, h.Logger)
			return
		}
	}

	encodeJSON(w, res, h.Logger)
}

// HandlePostEntry handles request to POST /entry
func (h *ContactHandler) HandlePostEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	entry := &services.Entry{}
	err := decodeAndValidate(r, entry)
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	res, err := h.Service.AddEntry(entry)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	encodeJSON(w, res, h.Logger)
}

// HandlePutEntry handles requests to PUT /entry
func (h *ContactHandler) HandlePutEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	entry := &services.Entry{}
	err := decodeAndValidate(r, entry)
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}
	if entry.ID == "" {
		Error(w, ErrNoID, http.StatusBadRequest, h.Logger)
		return
	}

	err = h.Service.UpdateEntry(entry)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}

	w.WriteHeader(http.StatusAccepted)
}

// HandleDeleteEntryByID handles requests to DELETE /entry/:id
func (h *ContactHandler) HandleDeleteEntryByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id == "" {
		Error(w, ErrNoIDParam, http.StatusBadRequest, h.Logger)
		return
	}

	err := h.Service.DeleteEntryByID(id)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleGetCSV handles requests to GET /csv/entry
func (h *ContactHandler) HandleGetCSV(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.Service.EntriesToCSV()
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	defer func() {
		res.Close()
		os.Remove(res.Name())
	}()

	w.Header().Set("Content-Disposition", "attachment; filename=entries.csv")
	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
	res.Seek(0, 0)
	io.Copy(w, res)
}

// HandlePostCSV handles request to POST /csv/entry
func (h *ContactHandler) HandlePostCSV(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, handle, err := r.FormFile("file")
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}
	defer file.Close()

	mimeType := handle.Header.Get("Content-Type")
	if mimeType != "text/csv" {
		Error(w, ErrInvalidFileType, http.StatusBadRequest, h.Logger)
		return
	}
	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmp.*.csv")
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	defer os.Remove(tmpFile.Name())
	_, err = io.Copy(tmpFile, file)

	res, err := h.Service.CSVToEntries(tmpFile)
	if err != nil {
		switch err {
		case services.ErrEmailExists:
			errRes := &postCSVError{
				Err:            err.Error(),
				InvalidEntries: res,
			}
			encodeJSON(w, errRes, h.Logger)
			return
		default:
			Error(w, err, http.StatusInternalServerError, h.Logger)
		}
	}
	w.WriteHeader(http.StatusAccepted)
}
