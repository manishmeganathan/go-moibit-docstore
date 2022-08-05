package docstore

import (
	"fmt"

	"github.com/manishmeganathan/go-moibit-sdk"
)

// DocStore represents an interface for working with MOIBit as Document Database.
// The DocStore treats all directories in the app root "/" as a Collection. Files in the root are ignored.
type DocStore struct {
	c *moibit.Client
}

// NewDocStore generates a new DocStore for a given moibit.Client.
func NewDocStore(client *moibit.Client) (*DocStore, error) {
	return &DocStore{client}, nil
}

// ListCollections returns a slice of Collection objects in the DocStore's root.
func (docstore *DocStore) ListCollections() ([]Collection, error) {
	// List files at the root "/"
	files, err := docstore.c.ListFiles("/")
	if err != nil {
		return nil, fmt.Errorf("failed to list files at root: %w", err)
	}

	// Declare collection accumulator and iterate over files
	collections := make([]Collection, 0, len(files))
	for _, file := range files {
		// If the file descriptor is a directory, split its paths and wrap
		// into a Collection while appending into the accumulator
		if file.IsDirectory {
			collections = append(collections, Collection{pathSplit(file.Directory)})
		}
	}

	// Return the collections
	return collections, nil
}

func (docstore *DocStore) GetCollection(name string) (Collection, error) {
	// File Status on the pathJoin("/", name)
	// Verify that file is a directory
	// If dir does not exist, create it if allowCreate is set
	// Wrap path into Collection and return

	return Collection{}, nil
}

func (docstore *DocStore) RemoveCollection(name string) bool {
	// Remove with RemoveDirectory() enabled for pathJoin("/", name)

	return false
}
