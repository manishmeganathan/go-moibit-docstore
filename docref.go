package docstore

import (
	"fmt"

	"github.com/manishmeganathan/go-moibit-sdk"
)

// DocRef represents a reference to a Document on MOIBit
type DocRef struct {
	path []string
	file moibit.FileDescriptor
}

func NewDocRef(file moibit.FileDescriptor) (DocRef, error) {
	// Fail if the file desc is for a directory
	if file.IsDirectory {
		return DocRef{}, fmt.Errorf("cannot create DocRef from directory")
	}

	return DocRef{pathSplit(file.Path), file}, nil
}

func (docref *DocRef) Get() (Document, error) {
	// Read file at docref.Path()
	// Pass Bytes into NewDocument
	// Return the Document

	return Document{}, nil
}

func (docref *DocRef) Set(doc Document) error {
	// Serialize doc.data
	// Write file at docref.Path()

	return nil
}

func (docref *DocRef) Remove() error {
	// Remove file at docref.Path()

	return nil
}

func (docref *DocRef) Exists() bool {
	return false
	//return docref.file.Exists()
}

func (docref *DocRef) Parent() Collection {
	return Collection{docref.path[:1]}
}

func (docref *DocRef) Path() string {
	return pathJoin(docref.path...)
}
