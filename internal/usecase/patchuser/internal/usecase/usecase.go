package usecase

import (
	"github.com/jennwah/ryde-backend-engineer/internal/usecase/patchuser/internal/contract"
)

type UseCase struct {
	repo contract.Repository
}

func New(repo contract.Repository) UseCase {
	return UseCase{
		repo: repo,
	}
}
