package task

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	mu   sync.Mutex
	path string
	//T in Tasks needs to be capital or else it cant be accessed, holy fcking shit!!
	Tasks []Task
}

func NewStore() (*Store, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(homedir, "Desktop", "Godo", "Tasks.json")
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)

	s := &Store{
		path:  path,
		Tasks: []Task{},
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s.Tasks)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s.Tasks)
}
func (s *Store) GenerateID() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	used := make(map[int]bool)
	for _, t := range s.Tasks {
		used[t.ID] = true
	}

	for {
		id := rand.Intn(9000) + 1000 // random 4-digit number
		if !used[id] {
			return id
		}
	}
}
