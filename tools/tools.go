package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

// The Data Models--------------------------------
type Person struct {
	Id        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
}

// The "Class"-----------------------------------
type Persons struct {
	Data           []Person
	dataFolderPath string
	dataFilename   string
	dataFilePath   string
}

// The optional "Class" Init----------
func NewPersons() *Persons {
	p := new(Persons)
	p.dataFolderPath = "./data"
	p.dataFilename = "persons.json"
	p.dataFilePath = path.Join(p.dataFolderPath, p.dataFilename)
	os.MkdirAll(p.dataFolderPath, 0644)
	_, chk := os.Stat(p.dataFilePath)
	if chk != nil {
		err := os.WriteFile(p.dataFilePath, []byte("[]"), 0644)
		if err != nil {
			panic("gehderned")
		}
	}
	byteValue, _ := os.ReadFile(p.dataFilePath)
	json.Unmarshal(byteValue, &p.Data)

	return p
}

// - The Methods----------------------------------
func (p *Persons) Reload() error {
	byteValue, _ := os.ReadFile(p.dataFilePath)
	err := json.Unmarshal(byteValue, &p.Data)
	return err
}

func (p *Persons) Save() error {
	dat, err := json.MarshalIndent(p.Data, "", "  ")
	if err != nil {
		return err
	}
	// fmt.Println(string(dat))
	os.WriteFile(p.dataFilePath, dat, 0644)
	return nil
}

// -------------------
func (p *Persons) Add(item Person) uuid.UUID {
	item.Id = uuid.New()
	p.Data = append(p.Data, item)
	p.Save()
	return item.Id
}

// -------------------
func (p *Persons) Get(id uuid.UUID) (Person, error) {
	for _, item := range p.Data {
		if item.Id == id {
			return item, nil
		}
	}
	return Person{}, fmt.Errorf("item with id %s not found", id.String())
}

// -------------------
func (p *Persons) Change(id uuid.UUID, person Person) (Person, error) {
	chk := false
	for idx, item := range p.Data {
		if item.Id == id {
			person.Id = id
			p.Data[idx] = person
			chk = true
		}
	}
	if !chk {
		return Person{}, fmt.Errorf("item with id %s not found", id.String())
	}
	p.Save()
	return person, nil
}

// -------------------
func (p *Persons) Delete(id uuid.UUID) (uuid.UUID, error) {
	var test []Person
	// var idx int
	for _, item := range p.Data {
		if item.Id != id {
			test = append(test, item)
		}
	}
	if len(p.Data) == len(test) {
		return id, fmt.Errorf("item with id %s not found", id.String())
	}
	fmt.Println(test)
	p.Data = test
	p.Save()

	return id, nil
}

// -------------------

//-------------------

//-------------------

// Some Tools-----------------------------------
