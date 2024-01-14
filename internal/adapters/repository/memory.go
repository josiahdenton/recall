package repository

import (
	"github.com/josiahdenton/recall/internal/domain"
)

type InMemoryStorage struct {
	tasks           []domain.Task
	cycles          []domain.Cycle
	accomplishments map[string]domain.Accomplishment
}

// TODO - add back again for debug testing...
