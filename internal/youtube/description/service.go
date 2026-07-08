package description

import (
	"context"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetDescription(ctx context.Context, description string) (Metadata, error) {
    return ProcessDescription(description), nil
}
