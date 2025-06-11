package models

func GetDefaultAdmin() *User {
	return &User{
		Id:       1,
		Username: "admin",
		Email:    "admin@admin.ru",
		Password: "admin",
		Admin: Bool{
			val: true,
			ok:  true,
		},
	}
}
