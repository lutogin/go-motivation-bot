package users

import (
	"context"
	"motivation-bot/logging"
	userDto "motivation-bot/users/dto"
)

type Service struct {
	repo   *Repo
	logger *logging.Logger
}

func NewService(repository *Repo, logger *logging.Logger) *Service {
	logger.Logger.Infoln("Registering service.")
	return &Service{repo: repository, logger: logger}
}

func (s *Service) Create(ctx context.Context, payload userDto.CreateUserDto) (string, error) {
	result, err := s.repo.Create(ctx, payload)
	return result, err
}

func (s *Service) GetById(ctx context.Context, payload userDto.GetUserByIdDto) (UserEntity, error) {
	result, err := s.repo.GetById(ctx, payload)
	return result, err
}

func (s *Service) GetAll(ctx context.Context, payload userDto.GetUsersDto) ([]UserEntity, error) {
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

func (s *Service) GetByAlertingDate(ctx context.Context, payload userDto.GetUserByAlertingDateDto) ([]UserEntity, error) {
	result, err := s.repo.GetUsersByAlertingDate(ctx, payload)
	return result, err
}
