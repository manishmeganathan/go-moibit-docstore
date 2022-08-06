package docstore

import (
	"fmt"

	"github.com/manishmeganathan/go-moibit-client"
)

// DocStore represents an interface for working with MOIBit as Document Database.
// The DocStore treats all directories in the app root "/" as a Collection. Files in the root are ignored.
type DocStore struct {
	client *moibit.Client
}

// NewDocStore generates a new DocStore for a given moibit.Client.
func NewDocStore(client *moibit.Client) (*DocStore, error) {
	return &DocStore{client}, nil
}

// ListCollections returns a slice of Collection objects in the DocStore's root.
func (docstore *DocStore) ListCollections() ([]*Collection, error) {
	// List files at the root "/"
	files, err := docstore.client.ListFiles("/")
	if err != nil {
		return nil, fmt.Errorf("failed to list files at root: %w", err)
	}

	// Declare collection accumulator and iterate over files
	collections := make([]*Collection, 0, len(files))
	for _, file := range files {
		// If the file descriptor is a directory, split its paths and wrap
		// into a Collection while appending into the accumulator
		if file.IsDirectory {
			collections = append(collections, &Collection{docstore.client, pathSplit(file.Directory)})
		}
	}

	// Return the collections
	return collections, nil
}

// GetCollection attempts to retrieve a Collection of the given name from the root of the DocStore.
// If the collection does not exist, it will be created.
func (docstore *DocStore) GetCollection(name string) (*Collection, error) {
	// Create Collection object with path to specified collection
	path := pathJoin(name)
	collection := &Collection{docstore.client, pathSplit(path)}

	// Attempt to make the directory at the specified path
	if err := docstore.client.MakeDirectory(path); err != nil {
		// If directory already exist, expect a 400 error with a specific message,
		// This indicated that the collection already exists and need not be created.
		if err.Error() == "non-ok response [400]: directory exist | directory already exist" {
			return collection, nil
		}

		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return collection, nil
}

// RemoveCollection removes a Collection from the DocStore.
// Calling it is idempotent and is a no-op if the collection does not exist
func (docstore *DocStore) RemoveCollection(name string) error {
	// Create path to the collection and attempt to remove it
	path := pathJoin(name)
	if err := docstore.client.RemoveFile(path, 0, moibit.RemoveDirectory()); err != nil {
		return fmt.Errorf("error removing collection: %w", err)
	}

	return nil
}
