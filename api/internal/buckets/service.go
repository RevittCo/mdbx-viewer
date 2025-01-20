package buckets

import (
	"errors"
	"fmt"
	"github.com/RevittConsulting/mdbx-viewer/internal/db_mdbx"
	"github.com/RevittConsulting/mdbx-viewer/types"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Db interface {
	// Open - Open the environment.
	Open(path string) error
	// Close - Close the environment.
	Close()
	// List - Returns all dbi names.
	List() ([]string, error)
	// Entries - Number of data items.
	Entries(name string) (uint64, error)
	// ValuesByKey - Returns all values with the given key.
	ValuesByKey(name string, key []byte) ([][]byte, error)
	// KeysByValue - Returns all keys with the given value.
	KeysByValue(name string, value []byte) ([][]byte, error)
	// Read - Read data from the database.
	Read(name string, take, offset uint64) ([]types.KeyValuePair, error)
	// CountKeysOfLength - Count and keys of a given length.
	CountKeysOfLength(name string, length uint64) (uint64, []string, error)
	// CountValuesOfLength - Count and values of a given length.
	CountValuesOfLength(bucketName string, length uint64) (uint64, []string, error)
}

type Service struct {
	db Db
	mu sync.RWMutex
}

func NewService(db Db) *Service {
	return &Service{
		db: db,
		mu: sync.RWMutex{},
	}
}

func (s *Service) GetDataSource() (*DataSource, error) {
	dir := os.Getenv("DATA_DIR")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	dataSource := &DataSource{Source: dir}
	for _, entry := range entries {
		treeElement, err := buildTreeElement(dir, entry)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		dataSource.TreeElements = append(dataSource.TreeElements, treeElement)
	}

	return dataSource, nil
}

func buildTreeElement(parentPath string, entry os.DirEntry) (TreeElement, error) {
	fullPath := filepath.Join(parentPath, entry.Name())
	treeElement := TreeElement{
		Name:         entry.Name(),
		FullPath:     fullPath,
		IsSelectable: true,
	}

	if entry.IsDir() {
		children, err := aggregateFolder(fullPath)
		if children == nil {
			return treeElement, nil
		}
		if err != nil {
			return TreeElement{}, err
		}
		treeElement.Children = children
	}

	return treeElement, nil
}

func aggregateFolder(path string) ([]TreeElement, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		if errors.Is(err, unix.EPERM) {
			return nil, nil
		}
		return nil, err
	}

	var children []TreeElement
	for _, entry := range entries {
		child, err := buildTreeElement(path, entry)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}

	return children, nil
}

func (s *Service) Open(path string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("database file does not exist at path: %s", path)
	}

	s.db.Close()

	newEnv := db_mdbx.New()
	if err := newEnv.Open(path); err != nil {
		return nil, err
	}

	s.db = newEnv

	return s.db.List()
}

func (s *Service) Read(bucketName string, number, length uint64) ([]types.KeyValuePairString, error) {
	foundData, err := s.db.Read(bucketName, length, number)
	if err != nil {
		return nil, err
	}

	data := make([]types.KeyValuePairString, 0)

	for _, kv := range foundData {
		data = append(data, kv.HexKeyHexValue())
	}

	return data, nil
}
