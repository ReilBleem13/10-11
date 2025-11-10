package repository

import "github.com/ReilBleem13/10-11/internal/dto"

type LinkRepository interface {
	Save(result dto.LinkCheckResult) error
	Load(id int) (dto.LinkCheckResult, error)
	NextID() (int, error)
}
