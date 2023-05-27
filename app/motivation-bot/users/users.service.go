package users

import (
	"context"
	"motivation-bot/logging"
	"motivation-bot/users/dto"
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

func (s *Service) GetAll(ctx context.Context, payload usersDto.GetUsersDto) ([]UserEntity, error) {
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

func (s *Service) GetByAlertingDate(ctx context.Context, payload usersDto.GetUserByAlertingTimeDto) ([]UserEntity, error) {
	from := payload.Date.Add(-5 * time.Minute)
	to := payload.Date.Add(5 * time.Minute)
	hoursFrom := from.Hour()
	minutesFrom := from.Minute()
	hoursTo := to.Hour()
	minutesTo := to.Minute()

	result, err := s.repo.GetUsersByAlertingTime(
		ctx,
		GetUsersByAlertingTimeDto{
			HoursFrom:   hoursFrom,
			MinutesFrom: minutesFrom,
			HoursTo:     hoursTo,
			MinutesTo:   minutesTo,
		},
	)

	return result, err
}
