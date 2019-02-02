package ukjent

// Word is a word + transation we are storing
type Word struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Note        string `json:"note"`
}
