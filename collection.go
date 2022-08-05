package docstore

import (
	"fmt"
	"strings"
)

// Collection represents a reference to a collection/directory on MOIBit
type Collection struct {
	path []string
}

func (collection *Collection) ListDocuments() ([]DocRef, error) {
	// ListFile at collection.Path()
	// Iterate through the FileDescriptors and create a DocRef for each of them if it is not a directory
	// Return the slice of created DocRefs

	return nil, nil
}

func (collection *Collection) ListCollections() ([]Collection, error) {
	// ListFile at collection.Path()
	// Iterate through the FileDescriptors and create a Collection for each of them if it is a directory
	// Return the slice of created Collections

	return nil, nil
}

func (collection *Collection) GetDocument(name string, allowCreate bool) (DocRef, error) {
	// FileStatus on the pathJoin(collection.path, name)
	// Verify file is not Directory
	// if file.Exist() == false
	// 	 allowCreate is true -> WriteFile at the above path
	//   else -> throw error
	// else wrap file descriptor into DocRef and return

	return DocRef{}, nil
}

func (collection *Collection) GetCollection(name string) (Collection, error) {
	// File Status on the pathJoin(collection.path, name)
	// Verify that file is a directory
	// If dir does not exist, create it if allowCreate is set
	// Wrap path into Collection and return

	return Collection{}, nil
}

func (collection *Collection) RemoveCollection(name string) error {
	// Remove with RemoveDirectory() enabled for pathJoin(collection.path, name)

	return nil
}

func (collection *Collection) Parent() (Collection, error) {
	if pathJoin(collection.path...) == "/" {
		return Collection{}, fmt.Errorf("collection has not parent")
	}

	return Collection{collection.path[:1]}, nil
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
