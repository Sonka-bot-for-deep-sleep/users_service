package mapper

import (
	pb "github.com/Sonka-bot-for-deep-sleep/proto_files/api"
	"github.com/Sonka-bot-for-deep-sleep/user_service/application/models"
)

func ToUser(user *models.User) *pb.User {
	return &pb.User{
		TgId:  user.TgId,
		Login: user.Login,
		Name:  user.Name,
	}
}
