package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
)

const MaxDisplay = 10

// Exercice 2
type tagCount struct {
	tag   string
	count int
}

// Exercice 3
type Note struct {
	ID      string
	Title   string
	Content string
	Tags    []string
}

func NewNote(title, content string) *Note {
	return &Note{Title: title, Content: content}
}

func (n *Note) Preview() string {
	lettres := []rune(n.Content)
	if len(lettres) > 80 {
		return string(lettres[:80])
	}
	return n.Content
}

func (n *Note) AddTag(tag string) {
	for _, t := range n.Tags {
		if t == tag {
			return
		}
	}
	n.Tags = append(n.Tags, tag)
}

func LoadFromFile(path string) ([]*Note, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var notes []*Note
	err = json.NewDecoder(f).Decode(&notes)
	return notes, err
}

// Exercice 4
var ErrDuplicate = errors.New("note already exists")
var ErrNotFound = errors.New("note not found")
var ErrValidation = errors.New("title required")

type NoteStore interface {
	Save(n *Note) error
	Get(id string) (*Note, error)
	All() []*Note
}

type MemoryStore struct {
	notes map[string]*Note
	mu    sync.Mutex
}

func (s *MemoryStore) Save(n *Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if n == nil || n.Title == "" {
		return ErrValidation
	}
	if _, exists := s.notes[n.ID]; exists {
		return ErrDuplicate
	}
	s.notes[n.ID] = n
	return nil
}

func (s *MemoryStore) Get(id string) (*Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, exists := s.notes[id]
	if !exists {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *MemoryStore) All() []*Note {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []*Note{}
	for _, n := range s.notes {
		result = append(result, n)
	}
	return result
}

func main() {
	// Exercice 1
	mots := os.Args[1:] // os.Args[0] est le nom du programme
	if len(mots) == 0 {
		fmt.Println("Erreur : aucun argument")
		os.Exit(1)
	}
	fmt.Println("--- Exercice 1 ---")
	fmt.Println("Nombre de mots :", len(mots))
	affiches := 0
	for _, mot := range mots {
		if len(mot) > 4 && affiches < MaxDisplay {
			fmt.Println(mot)
			affiches++
		}
	}
	var s string
	fmt.Printf("%q\n", s)
	x := 0
	var y int
	fmt.Println(x, y)

	// Exercice 2
	fmt.Println("--- Exercice 2 ---")
	occurrences := make(map[string]int)
	for _, mot := range mots {
		occurrences[mot]++
	}
	var tags []tagCount
	for mot, nombre := range occurrences {
		tags = append(tags, tagCount{tag: mot, count: nombre})
	}
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].count == tags[j].count {
			return tags[i].tag < tags[j].tag
		}
		return tags[i].count > tags[j].count
	})
	for _, t := range tags {
		if t.count > 1 {
			fmt.Println(t.tag, t.count)
		}
	}

	// Exercice 3
	fmt.Println("--- Exercice 3 ---")
	note := NewNote("Bonjour", "J'apprends Go avec des exercices simples")
	note.AddTag("go")
	note.AddTag("go") // Pas de doublon
	fmt.Println(note.Preview(), note.Tags)

	notes, err := LoadFromFile("notes.json")
	if err != nil {
		fmt.Println("Erreur lecture :", err)
	} else {
		for _, n := range notes {
			for _, tag := range n.Tags {
				if tag == "go" {
					fmt.Println("Note avec tag go :", n.Title)
					break
				}
			}
		}
	}

	// Exercice 4
	fmt.Println("--- Exercice 4 ---")
	store := &MemoryStore{notes: make(map[string]*Note)}
	var _ NoteStore = store
	note.ID = "1"
	fmt.Println("Save :", store.Save(note))
	fmt.Println("Doublon :", errors.Is(store.Save(note), ErrDuplicate))
	_, err = store.Get("inconnu")
	fmt.Println("Introuvable :", errors.Is(err, ErrNotFound))
	fmt.Println("Nombre de notes :", len(store.All()))
}
