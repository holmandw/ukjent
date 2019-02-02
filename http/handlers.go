package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/holman_dw/ukjent"
	"github.com/holman_dw/ukjent/json"
	"github.com/holman_dw/ukjent/store"

	"github.com/gorilla/mux"
)

var _ = http.MethodGet

func (s server) registerRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", s.Root)
	r.HandleFunc("/{word}", s.Get)
	return r
}

func (s server) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	word, ok := vars["word"]
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		errorResponse(w, "invalid url: expected `/{word}`", http.StatusBadRequest)
		return
	}

	entry, err := s.store.Get(word)
	if err != nil {
		switch err.(type) {
		case store.WordNotFoundError:
			errorResponse(w, err.Error(), http.StatusNotFound)
			return
		default:
			errorResponse(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if err := writeWord(w, entry); err != nil {
		log.Printf("error writing respose: %v", err)
	}
}

func (s server) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		writeWords(w, s.store.GetAll())
		return
	}
	if r.Method != http.MethodPost {
		errorResponse(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
	if c := r.Header.Get("Content-Type"); c != "application/json" {
		errorResponse(w, fmt.Sprintf("invalid Content-Type: %s", c), http.StatusBadRequest)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, "unable to read request body", http.StatusBadRequest)
		log.Printf("error reading body: %v", err)
		return
	}
	log.Print("got", string(b))
	word, err := json.UnmarshalWord(b)
	if err != nil {
		errorResponse(w, "unable to parse json body", http.StatusBadRequest)
		log.Printf("unable to parse json body: %v, %v", err, word)
		return
	}
	if err := s.store.Insert(word); err != nil {
		switch err.(type) {
		case store.TranslationExistsError:
			errorResponse(w, "word %v already exists", http.StatusConflict)
			log.Printf("cannot insert %v again", word)
			return
		default:
			errorResponse(w, "internal server error", http.StatusInternalServerError)
			log.Printf("error insertinv word %v %v", word, err)
			return
		}
	}
	if err := writeWord(w, word); err != nil {
		log.Printf("error sending response %v: %v", word, err)
	}
}

func (s server) Update(w http.ResponseWriter, r *http.Request) {}

func writeWord(w http.ResponseWriter, word ukjent.Word) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalWord(word)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = w.Write(b)
	return err
}
func writeWords(w http.ResponseWriter, words []ukjent.Word) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalWords(words)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = w.Write(b)
	return err
}

func errorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	e := ukjent.Error{
		Error: errorMsg,
		Code:  statusCode,
	}
	b, err := json.MarshalError(e)
	if err != nil {
		log.Printf("unhandled err marshalling error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(b)
}
