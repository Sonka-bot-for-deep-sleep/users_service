package models

type User struct {
	ID    int
	TgId  string `sql:"Tg_ID"`
	Name  string
	Login string
}
