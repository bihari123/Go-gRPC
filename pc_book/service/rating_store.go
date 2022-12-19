package service

import "sync"

type RatingStore interface {
	Add(laptopId string, score float64) (*Rating, error)
}

type Rating struct {
	Count uint32
	Sum   float64
}

type InMemoryRatingScore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

func NewInMemoryratingStore() *InMemoryRatingScore {
	return &InMemoryRatingScore{
		rating: make(map[string]*Rating),
	}
}

func (store *InMemoryRatingScore) Add(laptopId string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.rating[laptopId]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score
	}

	store.rating[laptopId] = rating

	return rating, nil
}
