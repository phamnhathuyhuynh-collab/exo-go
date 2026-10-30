package main

import (
	"errors"
	"testing"
)

// Test 1 : Save une note valide
func TestSave_valid(t *testing.T) {
	store := &MemoryStore{notes: make(map[string]*Note)}

	note := &Note{ID: "1", Title: "Go"}

	err := store.Save(note)

	if err != nil {
		t.Fatal(err)
	}
}

// Test 2 : Save une note sans titre
func TestSave_emptyTitle(t *testing.T) {
	store := &MemoryStore{notes: make(map[string]*Note)}

	note := &Note{ID: "1", Title: ""}

	err := store.Save(note)

	if !errors.Is(err, ErrValidation) {
		t.Fatal("Erreur attendue : ErrValidation")
	}
}

// Test 3 : Save deux fois la même note
func TestSave_duplicate(t *testing.T) {
	store := &MemoryStore{notes: make(map[string]*Note)}

	note := &Note{ID: "1", Title: "Go"}

	store.Save(note)
	err := store.Save(note)

	if !errors.Is(err, ErrDuplicate) {
		t.Fatal("Erreur attendue : ErrDuplicate")
	}
}

// Test 4 : Get une note inexistante
func TestGet_notFound(t *testing.T) {
	store := &MemoryStore{notes: make(map[string]*Note)}

	_, err := store.Get("?")

	if !errors.Is(err, ErrNotFound) {
		t.Fatal("Erreur attendue : ErrNotFound")
	}
}
