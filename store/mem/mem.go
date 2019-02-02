package mem

import (
	"sync"

	"github.com/holman_dw/ukjent"
	base "github.com/holman_dw/ukjent/store"
)

var (
	// Sanity check
	_ base.Store = &store{}
)

// Store is an in-memory version of a store
type store struct {
	mu   sync.RWMutex
	data map[string]entry
}

type entry struct {
	translation string
	note        string
}

// New builds a new Store
func New() base.Store {
	var s store
	s.data = make(map[string]entry)
	return &s
}

// Find a Word or return nil
func (m *store) Find(s string) *ukjent.Word {
	w, err := m.get(s)
	if err != nil {
		return nil
	}
	return &w
}

func (m *store) get(s string) (ukjent.Word, error) {
	defer m.withRead()()
	e, ok := m.data[s]
	if !ok {
		return empty(), base.NotFound(s)
	}
	w := word(s, e.translation, e.note)
	return w, nil
}

// Get a Word or return an error
func (m *store) Get(s string) (ukjent.Word, error) {
	return m.get(s)
}

func (m *store) GetAll() []ukjent.Word {
	words := make([]ukjent.Word, 0)
	for k, v := range m.data {
		words = append(words, word(k, v.translation, v.note))
	}
	return words
}

// Insert a word. Report an error if it already exists
func (m *store) Insert(w ukjent.Word) error {
	got, err := m.get(w.Word)
	if err == nil {
		return base.TranslationExists(got.Word, got.Translation)
	}

	defer m.withWrite()()
	m.data[w.Word] = entry{w.Translation, w.Note}
	return nil
}

// Update a Word. If it doesn't exist, Insert it.
func (m *store) Update(w ukjent.Word) error {
	defer m.withWrite()()
	m.data[w.Word] = entry{w.Translation, w.Note}
	return nil
}

func (m *store) withWrite() func() {
	m.mu.Lock()
	return m.mu.Unlock
}

func (m *store) withRead() func() {
	m.mu.RLock()
	return m.mu.RUnlock
}

func word(word string, transation string, note string) ukjent.Word {
	return ukjent.Word{
		Word:        word,
		Translation: transation,
		Note:        note,
	}
}

func empty() ukjent.Word {
	return ukjent.Word{}
}
