package main

import (
	"reflect"
	"testing"
)

func TestGetRelated(t *testing.T) {
	store := NewUserStore()
	user1 := "user1"
	user2 := "user2"
	user3 := "user3"
	store.Add(1, user1, []int{2, 3})
	store.Add(2, user2, []int{})
	store.Add(3, user3, []int{})

	relatedUsers := store.GetRelated(1)

	if totalRelatedUsers := len(relatedUsers); totalRelatedUsers != 2 {
		t.Errorf("Should only have 2 relations, not %d", totalRelatedUsers)
	}
	if !reflect.DeepEqual(relatedUsers[0], user2) {
		t.Errorf("Expected to find user 2 as a relation of user 1, got %d", relatedUsers[0])
	}
	if !reflect.DeepEqual(relatedUsers[1], user3) {
		t.Errorf("Expected to find user 3 as a relation of user 1, got %d", relatedUsers[1])
	}
}

func TestNoRecords(t *testing.T) {
	store := NewUserStore()

	relatedUsers := store.GetRelated(2)

	if totalRelatedUsers := len(relatedUsers); totalRelatedUsers != 0 {
		t.Errorf("Should not have any related users, got %d", totalRelatedUsers)
	}
}

func TestNotFoundRelationships(t *testing.T) {
	store := NewUserStore()
	user1 := "user1"
	store.Add(1, user1, []int{4, 5} /* User 4 and 5 don't have records */)

	relatedUsers := store.GetRelated(1)

	if totalRelatedUsers := len(relatedUsers); totalRelatedUsers != 0 {
		t.Errorf("User 1 shouldn't have any related users, got %d", totalRelatedUsers)
	}
}
