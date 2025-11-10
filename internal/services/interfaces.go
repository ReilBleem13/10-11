package services

import "github.com/ReilBleem13/10-11/internal/dto"

type File interface {
	Report(data dto.NewLinksNumRequest) ([]byte, error)
	Check(data dto.NewLinkRequest) ([]dto.LinkState, int, error)
}
