package docstore

import (
	"fmt"
	"strings"

	"github.com/manishmeganathan/go-moibit-sdk"
)

// Collection represents a reference to a collection/directory on MOIBit
type Collection struct {
	client *moibit.Client
	path   []string
}

// ListDocuments returns a slice of DocRef objects in the Collection.
func (collection *Collection) ListDocuments() ([]*DocRef, error) {
	// List files inside the collection path
	files, err := collection.client.ListFiles(collection.Path())
	if err != nil {
		return nil, fmt.Errorf("failed to list files in collection: %w", err)
	}

	// Declare document accumulator and iterate over the files
	documents := make([]*DocRef, 0, len(files))
	for _, file := range files {
		// If the file descriptor is a file, create a new DocRef
		// and append into the accumulator
		if !file.IsDirectory {
			doc, err := NewDocRef(file, collection.client)
			if err != nil {
				return nil, fmt.Errorf("failed to create docref for '%v': %w", file.Hash, err)
			}

			documents = append(documents, doc)
		}
	}

	// Return the documents
	return documents, nil
}

// ListCollections returns a slice of Collection
func (collection *Collection) ListCollections() ([]*Collection, error) {
	// List files inside the collection path
	files, err := collection.client.ListFiles(collection.Path())
	if err != nil {
		return nil, fmt.Errorf("failed to list files in collection: %w", err)
	}

	// Declare collection accumulator and iterate over files
	collections := make([]*Collection, 0, len(files))
	for _, file := range files {
		// If the file descriptor is a directory, split its paths and wrap
		// into a Collection while appending into the accumulator
		if file.IsDirectory {
			collections = append(collections, &Collection{collection.client, pathSplit(file.Directory)})
		}
	}

	return collections, nil
}

func (collection *Collection) GetDocument(name string, allowCreate bool) (*DocRef, error) {
	// FileStatus on the pathJoin(collection.path, name)
	// Verify file is not Directory
	// if file.Exist() == false
	// 	 allowCreate is true -> WriteFile at the above path
	//   else -> throw error
	// else wrap file descriptor into DocRef and return
	path := pathJoin(collection.Path(), name)

	fileDescriptor, err := collection.client.FileStatus(path)
	if err != nil {
		return nil, fmt.Errorf("error fetching file descriptor: %w", err)
	}

	if fileDescriptor.Exists() {
		return NewDocRef(fileDescriptor, collection.client)
	}

	if allowCreate {
		fileDescriptors, err := collection.client.WriteFile([]byte{}, path, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating new file: %w", err)
		}

		return NewDocRef(fileDescriptors[0], collection.client)
	}

	return nil, fmt.Errorf("error document not found")
}

func (collection *Collection) GetCollection(name string) (*Collection, error) {
	// Generate a new path from collection.path
	// Create a collection for the new path
	// Check if directory already exists by calling makeDirectory
	// Return the new collection
	path := pathJoin(collection.Path(), name)
	newCollection := &Collection{collection.client, pathSplit(path)}

	// Attempt to make the directory at the specified path
	if err := collection.client.MakeDirectory(path); err != nil {
		// If directory already exist, expect a 400 error with a specific message,
		// This indicated that the collection already exists and need not be created.
		if err.Error() == "non-ok response [400]: directory exist | directory already exist" {
			return newCollection, nil
		}

		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return newCollection, nil
}

func (collection *Collection) RemoveCollection(name string) error {
	// Remove with RemoveDirectory() enabled for pathJoin(collection.path, name)
	path := pathJoin(collection.Path(), name)

	if err := collection.client.RemoveFile(path, 0, moibit.RemoveDirectory()); err != nil {
		return fmt.Errorf("error removing collection: %w", err)
	}

	return nil
}

func (collection *Collection) Parent() (*Collection, error) {
	if pathJoin(collection.path...) == "/" {
		return nil, fmt.Errorf("collection has not parent")
	}

	return &Collection{collection.client, collection.path[:1]}, nil
}

func (collection *Collection) Path() string {
	return pathJoin(collection.path...)
}

func pathSplit(path string) []string {
	paths := make([]string, 0)
	for _, split := range strings.Split(path, "/") {
		if split != "" {
			paths = append(paths, split)
		}
	}

	return paths
}

func pathJoin(paths ...string) string {
	return fmt.Sprintf("/%v", strings.Join(paths, "/"))
}
