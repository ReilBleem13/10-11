package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ReilBleem13/10-11/internal/dto"
)

const dataDir = "data"

type fileRepository struct {
	mu      sync.Mutex
	nextNum int
}

func NewFileRepository() (LinkRepository, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	repo := &fileRepository{}
	if err := repo.loadNextID(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *fileRepository) loadNextID() error {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}

	maxID := 0
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		idStr := strings.TrimSuffix(f.Name(), ".json")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		if id > maxID {
			maxID = id
		}
	}

	r.mu.Lock()
	r.nextNum = maxID + 1
	if r.nextNum <= 1 {
		r.nextNum = 1
	}
	r.mu.Unlock()
	return nil
}

func (r *fileRepository) NextID() (int, error) {
	r.mu.Lock()
	id := r.nextNum
	r.nextNum++
	r.mu.Unlock()
	return id, nil
}

func (r *fileRepository) Save(result dto.LinkCheckResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	filename := filepath.Join(dataDir, fmt.Sprintf("%d.json", result.ID))
	return os.WriteFile(filename, data, 0644)
}

func (r *fileRepository) Load(id int) (dto.LinkCheckResult, error) {
	filename := filepath.Join(dataDir, fmt.Sprintf("%d.json", id))
	data, err := os.ReadFile(filename)
	if err != nil {
		return dto.LinkCheckResult{}, err
	}

	var result dto.LinkCheckResult
	if err := json.Unmarshal(data, &result); err != nil {
		return dto.LinkCheckResult{}, err
	}
	return result, nil
}
