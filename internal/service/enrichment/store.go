package enrichment

import "sync"

type Store struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewStore() *Store {
	return &Store{jobs: map[string]Job{}}
}

func (s *Store) Put(job Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

func (s *Store) Get(id string) (Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	return j, ok
}
