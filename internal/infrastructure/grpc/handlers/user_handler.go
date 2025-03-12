package handlers

import (
	"context"

	pb "github.com/Sonka-bot-for-deep-sleep/proto_files/api"
	"github.com/Sonka-bot-for-deep-sleep/user_service/application/dto"
	"github.com/Sonka-bot-for-deep-sleep/user_service/application/mapper"
	"github.com/Sonka-bot-for-deep-sleep/user_service/application/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service interface {
	GetByTgID(ctx context.Context, user *dto.GetUserByTgID) (*models.User, error)
	CreateUser(ctx context.Context, user *dto.CreateUser) error
}

type user struct {
	service service
	logger  *zap.Logger
	pb.UnimplementedUsersServiceServer
}

func New(service service, logger *zap.Logger) *user {
	return &user{
		service: service,
		logger:  logger,
	}
}

func (u *user) GetUserByTgID(
	ctx context.Context, in *pb.GetUserByTgIDRequest,
) (*pb.GetUserByTgIDResponse, error) {
	user, err := u.service.GetByTgID(ctx, &dto.GetUserByTgID{TgId: in.TgId})
	if err != nil {
		u.logger.Error("Failed get user by tg id")
		return nil, status.Error(codes.NotFound, "Не получилось найти ваши данные в системе")
	}
	return &pb.GetUserByTgIDResponse{
		User: mapper.ToUser(user),
	}, nil
}

func (u *user) CreateUser(ctx context.Context,
	in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := u.service.CreateUser(ctx,
		&dto.CreateUser{TgID: in.User.TgId, Name: in.User.Name, Login: in.User.Login},
	); err != nil {
		u.logger.Error("Failed create user")
		return nil, status.Error(codes.Internal,
			"Не получилось зарегистрировать вас в системе, попробуйте позже",
		)
	}

	return &pb.CreateUserResponse{}, nil
}
