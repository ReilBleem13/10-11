package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ReilBleem13/10-11/internal/dto"
	customerror "github.com/ReilBleem13/10-11/internal/errors"
	"github.com/ReilBleem13/10-11/internal/services"
)

type Handler struct {
	srv services.File
}

func NewHandler(srv services.File) (Handler, error) {
	return Handler{
		srv: srv,
	}, nil
}

func (h *Handler) HandleCheck(w http.ResponseWriter, r *http.Request) {
	var newRequest dto.NewLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&newRequest); err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}

	result, id, err := h.srv.Check(newRequest)
	if err != nil {
		http.Error(w, "failed to check", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"links":     result,
		"links_num": id,
	}); err != nil {
		http.Error(w, "encode error", http.StatusBadRequest)
		return
	}
}

func (h *Handler) HandleReport(w http.ResponseWriter, r *http.Request) {
	var newRequest dto.NewLinksNumRequest
	if err := json.NewDecoder(r.Body).Decode(&newRequest); err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}

	if len(newRequest.LinksList) == 0 {
		http.Error(w, "links_list is empty", 400)
		return
	}

	pdfBytes, err := h.srv.Report(newRequest)
	if err != nil {
		if errors.Is(err, customerror.NoDataFound) {
			http.Error(w, customerror.NoDataFound.Error(), 404)
			return
		}
		http.Error(w, "failed to report", 500)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"report.pdf\"")
	w.Header().Set("Content-Length", fmt.Sprint(len(pdfBytes)))
	w.Write(pdfBytes)
}
