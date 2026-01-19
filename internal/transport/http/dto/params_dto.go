package dto

// IDParams para validar parámetros de ID en la URL
type IDParams struct {
	ID uint `uri:"id" validate:"required,min=1"`
}

// DNIParams para validar parámetros de DNI en la URL
type DNIParams struct {
	Dni string `uri:"dni" validate:"required,len=8,numeric"`
}

// ListUsersQuery para validar query parameters en el listado
type ListUsersQuery struct {
	Page  int    `form:"page" validate:"omitempty,min=1"`
	Limit int    `form:"limit" validate:"omitempty,min=1,max=100"`
	Email string `form:"email" validate:"omitempty,email"`
	Role  string `form:"role" validate:"omitempty,alpha"`
}
