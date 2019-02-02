package store

import (
	"fmt"

	"github.com/holman_dw/ukjent"
)

// Store is a store, yo
type Store interface {
	// Find a Word or return nil
	Find(string) *ukjent.Word
	// Get a Word or return an error
	Get(string) (ukjent.Word, error)
	// GetAll gets all words currently stored.
	GetAll() []ukjent.Word
	// Insert a word. Report an error if it already exists
	Insert(ukjent.Word) error
	// Update a Word. If it doesn't exist, Insert it.
	Update(ukjent.Word) error
}

// WordNotFoundError is used when we cannot get a word from a store
type WordNotFoundError struct {
	Word string
}

func (e WordNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Word)
}

// NotFound is a convenience
func NotFound(w string) WordNotFoundError {
	return WordNotFoundError{Word: w}
}

// TranslationExistsError is used when a translation exists
type TranslationExistsError struct {
	Word       string
	Transation string
}

func (e TranslationExistsError) Error() string {
	return fmt.Sprintf("%s - %s grouping already exists", e.Word, e.Transation)
}

// TranslationExists is a convenience
func TranslationExists(w string, t string) TranslationExistsError {
	return TranslationExistsError{Word: w, Transation: t}
}
