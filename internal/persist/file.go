package persist

import (
	"encoding/json"
	"os"
	"path/filepath"

	"example.com/peerdisco/internal/member"
)

type Snapshot struct {
	SelfID   string                       `json:"self_id"`
	SelfAddr string                       `json:"self_addr"`
	Members  []member.Member              `json:"members"`
	Meta     map[string]map[string]string `json:"meta"`
}

type Store struct {
	path string
}

func New(path string) *Store { return &Store{path: path} }

func (s *Store) Save(snap Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Store) Load() (Snapshot, error) {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
