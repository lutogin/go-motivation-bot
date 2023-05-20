package users

import (
	"context"
	"motivation-bot/common/logging"
	userDto "motivation-bot/users/dto"
)

type Service struct {
	repo   Repository
	logger *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) *Service {
	logger.Logger.Infoln("Registering service.")
	return &Service{repo: repository, logger: logger}
}

func (s *Service) Create(ctx context.Context, payload userDto.CreateUserDto) (id string, err error) {
	result, err := s.repo.Create(ctx, payload)
	return result, err
}

func (s *Service) GetById(ctx context.Context, payload userDto.GetUserByIdDto) (user UserEntity, err error) {
	result, err := s.repo.GetById(ctx, payload)
	return result, err
}

func (s *Service) GetAll(ctx context.Context, payload userDto.GetUsersDto) (users []UserEntity, err error) {
	result, err := s.repo.GetByFilter(ctx, payload)
	return result, err
}

func (s *Service) Update(ctx context.Context, payload userDto.UpdateUserDto) error {
	err := s.repo.Update(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, payload userDto.DeleteUserDto) error {
	err := s.repo.Delete(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}
