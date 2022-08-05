package docstore

import (
	"encoding/json"

	"github.com/manishmeganathan/go-moibit-sdk"
)

// Document represents a document on the MOIBit DocStore
type Document struct {
	name string
	data jsonIR
	desc moibit.FileDescriptor
}

// NewDocument returns a new Document for a given document name and some JSON encoded data.
// Returns an error if the data is not JSON encoded.
func NewDocument(name string, data []byte) (*Document, error) {
	// Create a new IR JSON and deserialize the data into it
	ir := make(jsonIR)
	if err := ir.Deserialize(data); err != nil {
		return nil, err
	}

	// Wrap the IR JSON in a Document and return
	return &Document{name, ir, moibit.FileDescriptor{}}, nil
}

func (doc *Document) Name() string {
	return doc.name
}

func (doc *Document) SetKey(key string, value any) {
	doc.data[key] = value
}

func (doc *Document) GetKey(key string) any {
	return doc.data[key]
}

type jsonIR map[string]any

func (ir *jsonIR) Serialize() ([]byte, error) {
	return json.Marshal(ir)
}

func (ir *jsonIR) Deserialize(data []byte) error {
	object := make(jsonIR)

	if len(data) != 0 {
		if err := json.Unmarshal(data, object); err != nil {
			return err
		}
	}

	*ir = object
	return nil
}
