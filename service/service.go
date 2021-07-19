package service

import (
	"context"
	"github.com/amrHassanAbdallah/notificationaway/persistence"
)

type Service struct {
	persistence persistence.PersistenceLayerInterface
}

// NewService returns new advisor manager that allows CRUD operations
func NewService(persistence persistence.PersistenceLayerInterface) *Service {
	return &Service{
		persistence: persistence,
	}
}
func (s Service) AddMessage(ctx context.Context, m persistence.Message) (*persistence.Message, error) {
	return s.persistence.AddMessage(ctx, m)
}

func (s Service) GetMessage(ctx context.Context, Id string) (*persistence.Message, error) {
	return s.persistence.GetMessage(ctx, Id)
}
