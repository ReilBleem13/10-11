package services

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/ReilBleem13/10-11/internal/dto"
	customerrors "github.com/ReilBleem13/10-11/internal/errors"
	"github.com/ReilBleem13/10-11/internal/repository"
	"github.com/jung-kurt/gofpdf"
)

type fileServices struct {
	repo repository.LinkRepository
}

func NewFileServices(repo repository.LinkRepository) (File, error) {
	return &fileServices{
		repo: repo,
	}, nil
}

func (f *fileServices) Report(data dto.NewLinksNumRequest) ([]byte, error) {
	var report strings.Builder
	found := false

	for _, id := range data.LinksList {
		s, err := f.repo.Load(id)
		if err != nil {
			continue
		}

		found = true
		report.WriteString(fmt.Sprintf("Report set #%d\n", id))
		for _, res := range s.Results {
			report.WriteString(fmt.Sprintf("%s : %s\n", res.URL, res.Status))
		}
		report.WriteString("\n")
	}

	if !found {
		return nil, customerrors.NoDataFound
	}

	pdfBytes, err := generatePDF(report.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF")
	}
	return pdfBytes, nil
}

func (f *fileServices) Check(data dto.NewLinkRequest) ([]dto.LinkState, int, error) {
	result := make([]dto.LinkState, 0, len(data.Links))
	for _, link := range data.Links {
		status := checkLink(link)
		result = append(result, dto.LinkState{URL: link, Status: status})
	}

	id, err := f.repo.NextID()
	if err != nil {
		return nil, 0, err
	}

	if err := f.repo.Save(dto.LinkCheckResult{ID: id, Results: result}); err != nil {
		return nil, 0, err
	}
	return result, id, nil
}

func checkLink(url string) string {
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return "not avaliable"
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "available"
	}
	return "not avaliable"
}

func generatePDF(text string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Helvetica", "", 12)
	pdf.MultiCell(0, 8, text, "", "", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
