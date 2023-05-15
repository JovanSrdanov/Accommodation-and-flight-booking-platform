package model

type UpdateProfileDto struct {
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Email   string  `json:"email"`
	Address Address `json:"address"`
}

func NewUpdateProfileDto(name, surname, email string, address Address) *UpdateProfileDto {
	return &UpdateProfileDto{
		Name:    name,
		Surname: surname,
		Email:   email,
		Address: address,
	}
}
