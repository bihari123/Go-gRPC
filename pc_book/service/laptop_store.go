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
	Search(filter *pcbook.Filter, found func(laptop *pcbook.Laptop) error) error
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

func deepCopy(laptop *pcbook.Laptop) (*pcbook.Laptop, error) {
	other := &pcbook.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("can't copy the laptop data: %v", err)
	}
	return other, err
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

func (store *InMemoryLaptopStore) Search(
	filter *pcbook.Filter,
	found func(laptop *pcbook.Laptop) error,
) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for _, laptop := range store.data {
		if IsQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func IsQualified(filter *pcbook.Filter, laptop *pcbook.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pcbook.Memory) uint64 {
	val := memory.GetValue()

	switch memory.GetUnit() {
	case pcbook.Memory_BIT:
		return val
	case pcbook.Memory_BYTE:
		return val << 3 // means val *8 (8 = 2^3)
	case pcbook.Memory_KILOBYTE:
		return val << 13 // 1024*8= 2^10 * 2^3 = 2^13
	case pcbook.Memory_MEGABYTE:
		return val << 23
	case pcbook.Memory_GIGABYTE:
		return val << 33
	case pcbook.Memory_TERABYTE:
		return val << 43
	default:
		return 0
	}
}
