package dto

type CreateUser struct {
	TgID  string `json:"tg_id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}
