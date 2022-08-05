package docstore

import (
	"encoding/json"
)

// Document represents a document on the MOIBit DocStore
type Document map[string]any

// NewDocument returns a new Document for some JSON encoded data.
// Returns an error if the data is not JSON encoded.
func NewDocument(data []byte) (Document, error) {
	// Create a new Document and set its data
	doc := make(Document)
	if err := doc.SetJSON(data); err != nil {
		return nil, err
	}

	return doc, nil
}

// SetKey sets a key-value pair into the Document
func (doc Document) SetKey(key string, value any) {
	doc[key] = value
}

// GetKey returns a value for a given key from the Document.
// Returns nil if key does not exist in the Document
func (doc Document) GetKey(key string) any {
	return doc[key]
}

// SetJSON accepts JSON bytes and overwrites the existing Document data
func (doc *Document) SetJSON(data []byte) error {
	// Create a new Document and attempt to deserialize the data into it
	// The data will not be unmarshalled if it is not JSON encoded
	newdoc := make(Document)
	if len(data) != 0 {
		// Decoding is skipped if data has no contents
		if err := json.Unmarshal(data, &newdoc); err != nil {
			return err
		}
	}

	// Set newdoc to doc
	*doc = newdoc
	return nil
}

// GetJSON returns the JSON Bytes representing the Document data
func (doc Document) GetJSON() ([]byte, error) {
	return json.Marshal(doc)
}
