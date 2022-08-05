package docstore

import "github.com/manishmeganathan/go-moibit-sdk"

type DocStore struct {
	c *moibit.Client
}

func NewDocStore(client *moibit.Client) (*DocStore, error) {
	return nil, nil
}

func (docstore *DocStore) GetCollection(name string, allowCreate bool) (Collection, error) {
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
