package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"addressbook.com/basicRest/contacts"
)

type contactServer struct {
	store *contacts.AddressBook
}

func (ab *contactServer) contactHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/contact/" {

		if req.Method == http.MethodPost {
			ab.createContactHandler(w, req)
		} else if req.Method == http.MethodGet {
			ab.getAllContactsHandler(w, req)
		} else if req.Method == http.MethodDelete {
			ab.deleteAllContactsHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /contact/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {

		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /contact/<emailid> in contact handler", http.StatusBadRequest)
			return
		}
		emailid := pathParts[1]

		if req.Method == http.MethodDelete {
			ab.deleteContactHandler(w, req, emailid)
		} else if req.Method == http.MethodGet {
			ab.getContactHandler(w, req, emailid)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /contact/<emailid>, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (cs *contactServer) createContactHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling contact create at %s\n", req.URL.Path)

	//Used Internally to converto to/from JSON
	type Request struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Address   string    `json:"address"`
		Mobile    string    `json:"mobile"`
		CreatedAt time.Time `json:"createdAt"`
	}

	type Response struct {
		Email string `json:"email"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rs Request
	if err := dec.Decode(&rs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := cs.store.CreateContact(rs.Name, rs.Email, rs.Address, rs.Mobile, time.Now())
	js, err := json.Marshal(Response{Email: email})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (cs *contactServer) getContactHandler(w http.ResponseWriter, req *http.Request, emailid string) {
	log.Printf("handling get contact at %s\n", req.URL.Path)

	contact, err := cs.store.GetContact(emailid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func (cs *contactServer) getAllContactsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Getting All Contacts at %s\n", req.URL.Path)

	allContacts := cs.store.GetAllContacts()
	js, err := json.Marshal(allContacts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (cs *contactServer) deleteContactHandler(w http.ResponseWriter, req *http.Request, emailid string) {
	log.Printf("handling delete contact at %s\n", req.URL.Path)

	err := cs.store.DeleteContact(emailid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (cs *contactServer) deleteAllContactsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all contact at %s\n", req.URL.Path)
	cs.store.DeleteAllContacts()
}

func (cs *contactServer) createdAtHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by due at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /due/<date>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}
	if len(pathParts) != 4 {
		badRequestError()
		return
	}

	year, err := strconv.Atoi(pathParts[1])
	if err != nil {
		badRequestError()
		return
	}
	month, err := strconv.Atoi(pathParts[2])
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, err := strconv.Atoi(pathParts[3])
	if err != nil {
		badRequestError()
		return
	}

	contacts := cs.store.GetContactByCreatedDate(year, time.Month(month), day)
	js, err := json.Marshal(contacts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func NewContactServer() *contactServer {
	store := contacts.New()
	return &contactServer{store: store}
}

func main() {
	mux := http.NewServeMux()
	server := NewContactServer()
	mux.HandleFunc("/contact/", server.contactHandler)
	mux.HandleFunc("/createdAt/", server.createdAtHandler)
	log.Fatal((http.ListenAndServe("localhost:8000", mux)))
}
