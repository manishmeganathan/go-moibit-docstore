package docstore

import (
	"fmt"
	"strings"

	"github.com/manishmeganathan/go-moibit-client"
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
			doc, err := newDocRef(file, collection.client)
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

// GetDocument returns a DocRef to the document in the collection with the given name.
// If the document does not exist, an error is thrown unless allowCreate is set, in which
// case, a blank document is created and stored before a reference to it is returned.
func (collection *Collection) GetDocument(name string, allowCreate bool) (*DocRef, error) {
	// Create path to the document and get its file stat
	path := pathJoin(collection.Path(), extAdd(name))
	file, err := collection.client.FileStatus(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: filestat check failed %w", err)
	}

	if !file.Exists() {
		// If the file does not exist and allowCreate is not set, return error
		if !allowCreate {
			return nil, fmt.Errorf("failed to get document: does not exist")
		}

		// allowCreate is set, so create a blank file at the path
		file, err := collection.client.WriteFile([]byte{}, path)
		if err != nil {
			return nil, fmt.Errorf("error creating new file: %w", err)
		}

		// Create a DocRef for the newly written file and return it
		return newDocRef(file, collection.client)
	}

	// Create a new DocRef for the existing file and return it
	return newDocRef(file, collection.client)
}

// GetCollection attempts to retrieve a Collection of the given name from within the calling collection.
// If the collection does not exist, it will be created.
func (collection *Collection) GetCollection(name string) (*Collection, error) {
	// Create Collection object with path to specified collection
	path := pathJoin(collection.Path(), name)
	subcol := &Collection{collection.client, pathSplit(path)}

	// Attempt to make the directory at the specified path
	if err := collection.client.MakeDirectory(path); err != nil {
		// If directory already exist, expect a 400 error with a specific message,
		// This indicated that the collection already exists and need not be created.
		if err.Error() == "non-ok response [400]: directory exist | directory already exist" {
			return subcol, nil
		}

		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return subcol, nil
}

// RemoveCollection removes a Collection from within the calling Collection.
// Calling it is idempotent and is a no-op if the collection does not exist.
func (collection *Collection) RemoveCollection(name string) error {
	// Create path to the sub-collection and attempt to remove it
	path := pathJoin(collection.Path(), name)
	if err := collection.client.RemoveFile(path, 0, moibit.RemoveDirectory()); err != nil {
		return fmt.Errorf("error removing collection: %w", err)
	}

	return nil
}

// Parent returns the parent Collection for the collection.
// Returns an error if the collection's parent would be the root.
// Note: Use the DocStore for accessing other collections in the root.
func (collection *Collection) Parent() (*Collection, error) {
	// If the path elements of the collection is 1, i.e, only
	// contains the collection name with no parent, return an error
	if len(collection.path) == 1 {
		return nil, fmt.Errorf("collection parent is the docstore root")
	}

	// Create a collection with the last path element removed
	return &Collection{collection.client, collection.path[:1]}, nil
}

// Path returns a path to the collection directory
func (collection *Collection) Path() string {
	return pathJoin(collection.path...)
}

// pathSplit is a utility function for splitting a path into its elements
func pathSplit(path string) []string {
	paths := make([]string, 0)
	for _, split := range strings.Split(path, "/") {
		if split != "" {
			paths = append(paths, split)
		}
	}

	return paths
}

// pathJoin is a utility function for joining path elements into a string
func pathJoin(paths ...string) string {
	return fmt.Sprintf("/%v", strings.Join(paths, "/"))
}

// extAdd is a utility function for adding the '.json' extension to a file name
func extAdd(file string) string {
	return strings.Join([]string{file, "json"}, ".")
}

// extRemove is a utility function for removing the '.json' extension to a file name
func extRemove(filename string) string {
	return strings.TrimPrefix(strings.Split(filename, ".json")[0], "/")
}
