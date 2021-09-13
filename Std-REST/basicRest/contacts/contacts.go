package contacts

import (
	"fmt"
	"sync"
	"time"
)

type Contacts struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Mobile    string    `json:"mobile"`
	CreatedAt time.Time `json:"createdAt"`
}

// Creating an In memory Address book that is safe to access concurrently
type AddressBook struct {
	sync.Mutex

	contacts map[string]Contacts
}

func New() *AddressBook {
	ab := &AddressBook{}
	ab.contacts = make(map[string]Contacts)

	return ab
}

// Create a new Contact in the address book
func (ab *AddressBook) CreateContact(name string, email string, address string, mobile string, createdAt time.Time) string {
	ab.Lock()

	defer ab.Unlock()

	contacts := Contacts{

		Name:      name,
		Email:     email,
		Address:   address,
		Mobile:    mobile,
		CreatedAt: createdAt,
	}
	ab.contacts[email] = contacts

	return email

}

// Retrieve a Contact by Emailid from the store and return an Error if it does
//exist
func (ab *AddressBook) GetContact(email string) (Contacts, error) {
	ab.Lock()

	defer ab.Unlock()
	c, ok := ab.contacts[email]

	if ok {
		return c, nil
	} else {
		return Contacts{}, fmt.Errorf("Contact with email=%s not found", email)
	}
}

//Delete a Contact with the give email,returns error if no match if found
func (ab *AddressBook) DeleteContact(email string) error {
	ab.Lock()
	defer ab.Unlock()

	if _, ok := ab.contacts[email]; !ok {
		return fmt.Errorf("Contact with Email=%s not found", email)
	}
	delete(ab.contacts, email)
	return nil
}

//Delete all Contacts
func (ab *AddressBook) DeleteAllContacts() error {
	ab.Lock()

	defer ab.Unlock()

	ab.contacts = make(map[string]Contacts)

	return nil
}

//Retrieve all Contacts
func (ab *AddressBook) GetAllContacts() []Contacts {
	ab.Lock()

	defer ab.Unlock()

	contacts := make([]Contacts, len(ab.contacts))
	for _, contact := range ab.contacts {
		contacts = append(contacts, contact)
	}
	fmt.Println(contacts)
	return contacts
}

//Get Contacts by CreatedAt Date
func (ab *AddressBook) GetContactByCreatedDate(year int, month time.Month, day int) []Contacts {
	ab.Lock()

	defer ab.Unlock()

	var contacts []Contacts

	for _, contact := range ab.contacts {
		y, m, d := contact.CreatedAt.Date()
		if y == year && m == month && d == day {
			contacts = append(contacts, contact)
		}
	}
	return contacts
}
