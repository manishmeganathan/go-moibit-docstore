package docstore

import (
	"fmt"

	"github.com/manishmeganathan/go-moibit-sdk"
)

// DocRef represents a reference to a Document on MOIBit
type DocRef struct {
	path   []string
	client *moibit.Client
	file   moibit.FileDescriptor
}

// newDocRef generates and returns a new DocRef instance for a given MOIBit File Descriptor and Client
func newDocRef(file moibit.FileDescriptor, client *moibit.Client) (*DocRef, error) {
	// Fail if the file desc is for a directory
	if file.IsDirectory {
		return nil, fmt.Errorf("cannot create DocRef from directory")
	}

	// Create a DocRef with the path to the document including its name and directory
	return &DocRef{pathSplit(pathJoin(file.Directory, file.Path)), client, file}, nil
}

// Get attempts to retrieve the Document at the DocRef
func (docref *DocRef) Get() (Document, error) {
	// ReadFile at the path of DocRef with version specified by the FileDescriptor
	data, err := docref.client.ReadFile(docref.Path(), docref.file.Version)
	if err != nil {
		return nil, err
	}

	// Create a new Document with the data from the file
	doc, err := NewDocument(data)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// Set attempts to set the given Document as the data for DocRef
func (docref *DocRef) Set(doc Document) error {
	// Get the JSON Bytes for the Document
	data, err := doc.GetJSON()
	if err != nil {
		return err
	}

	// WriteFile to path specified in DocRef
	if _, err = docref.client.WriteFile(data, docref.Path()); err != nil {
		return err
	}

	return nil
}

// Remove attempts to remove the Document at DocRef
func (docref *DocRef) Remove() error {
	// RemoveFile at path specified by DocRef with version from the FileDescriptor
	if err := docref.client.RemoveFile(docref.Path(), docref.file.Version); err != nil {
		return err
	}

	return nil
}

// Exists returns whether a Document exists at DocRef
func (docref *DocRef) Exists() bool {
	return docref.file.Exists()
}

// Parent returns the parent Collection of Document at DocRef
func (docref *DocRef) Parent() *Collection {
	return &Collection{docref.client, docref.path[:1]}
}

// Path returns the path to the Document at DocRef
func (docref *DocRef) Path() string {
	return pathJoin(docref.path...)
}
