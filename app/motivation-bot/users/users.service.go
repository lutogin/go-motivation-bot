package users

import (
	"context"
	"fmt"
	"motivation-bot/logging"
	"motivation-bot/users/dto"
	"strconv"
	"time"
)

type Service struct {
	repo   *Repo
	logger *logging.Logger
}

func NewService(repository *Repo, logger *logging.Logger) *Service {
	logger.Logger.Infoln("Registering service.")
	return &Service{repo: repository, logger: logger}
}

func (s *Service) Create(ctx context.Context, payload usersDto.CreateUserDto) (string, error) {
	result, err := s.repo.Create(ctx, payload)

	return result, err
}

func (s *Service) GetById(ctx context.Context, payload usersDto.GetUserByIdDto) (UserEntity, error) {
	result, err := s.repo.GetById(ctx, payload)

	return result, err
}

func (s *Service) GetByFilter(ctx context.Context, payload usersDto.GetUsersDto) ([]UserEntity, error) {
	result, err := s.repo.GetByFilter(ctx, payload)
	return result, err
}

func (s *Service) Update(ctx context.Context, payload usersDto.UpdateUserDto) error {
	err := s.repo.Update(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, payload usersDto.DeleteUserDto) error {
	err := s.repo.Delete(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteByUserName(ctx context.Context, payload usersDto.DeleteUserByUserNameDto) error {
	err := s.repo.DeleteByUserName(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteByChatId(ctx context.Context, payload usersDto.DeleteUserByChatIdDto) error {
	err := s.repo.DeleteByChatId(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetByAlertingDate(ctx context.Context, payload usersDto.GetUserByAlertingTimeDto) ([]UserEntity, error) {
	from := payload.Date.Add(-5 * time.Minute)
	to := payload.Date.Add(5 * time.Minute)
	timeFrom, _ := strconv.Atoi(fmt.Sprintf("%d%d", from.Hour(), from.Minute()))
	timeTo, _ := strconv.Atoi(fmt.Sprintf("%d%d", to.Hour(), to.Minute()))

	result, err := s.repo.GetUsersByAlertingTime(
		ctx,
		GetUsersByAlertingTimeDto{
			TimeFrom: timeFrom,
			TimeTo:   timeTo,
		},
	)

	return result, err
}

func (s *Service) Upsert(ctx context.Context, payload usersDto.UpdateUserDto) error {
	err := s.repo.Upsert(ctx, payload)

	return err
}

func (s *Service) GetByChatId(ctx context.Context, payload usersDto.GetUserByChatIdDto) (UserEntity, error) {
	user, err := s.repo.GetByChatId(ctx, payload)
	if err != nil {
		return UserEntity{}, err
	}

	return user, nil
}
