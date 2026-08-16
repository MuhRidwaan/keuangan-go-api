package service

import (
	"net/http"

	"keuangan-api/app/model"
	"keuangan-api/app/repository"
)

type APIDocService struct {
	Repo *repository.APIDocRepository
}

func (s *APIDocService) GetAll() ([]model.APIDocumentation, int, error) {
	docs, err := s.Repo.GetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return docs, http.StatusOK, nil
}
