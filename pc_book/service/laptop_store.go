package service

import (
	"fmt"
	pcbook "pcbook/proto"
	"sync"

	"github.com/jinzhu/copier"
)

type LaptopStore interface {
	Save(laptop *pcbook.Laptop) error
	Find(id string) (*pcbook.Laptop, error)
}

type InMemoryLaptopStore struct {
	mutex sync.Mutex
	data  map[string]*pcbook.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pcbook.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *pcbook.Laptop) error {
	store.mutex.Lock()

	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return fmt.Errorf("Laptop with id: %s already exists", laptop.Id)
	}

	// deep copy
	other := &pcbook.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("can't copy the laptop data: %v", err)
	}

	store.data[other.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*pcbook.Laptop, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	laptop := store.data[id]

	if laptop == nil {
		return nil, nil
	}

	other := &pcbook.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("error while deep copying :%w", err)
	}

	return other, err
}
