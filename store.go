package main

import (
	"sync"
)

type Store interface {
	Add(id int, data interface{}, relationships []int)
	GetRelated(id int) []interface{}
	Remove(id int)
}

type UserStore struct {
	users *sync.Map
}

type record struct {
	id            int
	data          interface{}
	relationships []int
}

func (s *UserStore) Add(id int, data interface{}, relationships []int) {
	s.users.Store(id, record{
		id,
		data,
		relationships,
	})
}

func (s *UserStore) GetRelated(id int) (relatedRecords []interface{}) {
	if user, ok := s.users.Load(id); ok {
		for _, relationshipID := range user.(record).relationships {
			if relatedRecord, ok := s.users.Load(relationshipID); ok {
				relatedRecords = append(relatedRecords, relatedRecord.(record).data)
			}
		}
	}

	return
}

func (s *UserStore) Remove(id int) {
	s.users.Delete(id)
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: &sync.Map{},
	}
}
