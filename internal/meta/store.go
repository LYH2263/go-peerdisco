package meta

// Store 节点元数据。
type Store struct {
	byID map[string]map[string]string
}

func NewStore() *Store {
	return &Store{byID: make(map[string]map[string]string)}
}

func (s *Store) Set(id string, kv map[string]string) {
	s.byID[id] = CloneMap(kv)
}

func (s *Store) Get(id string) map[string]string {
	return s.byID[id]
}

func (s *Store) GetClone(id string) map[string]string {
	return CloneMap(s.byID[id])
}

func (s *Store) Delete(id string) {
	delete(s.byID, id)
}

func (s *Store) AllClone() map[string]map[string]string {
	out := make(map[string]map[string]string, len(s.byID))
	for k, v := range s.byID {
		out[k] = CloneMap(v)
	}
	return out
}

func (s *Store) Replace(in map[string]map[string]string) {
	s.byID = make(map[string]map[string]string, len(in))
	for k, v := range in {
		s.byID[k] = CloneMap(v)
	}
}

func (s *Store) Clear() {
	s.byID = make(map[string]map[string]string)
}
