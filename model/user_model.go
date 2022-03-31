package model

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" `
	Password string `json:"password"`
	Role     string `json:"role"`
	Deleted  bool   `json:"deleted"`
}

type CreateStaff struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" validate:"password"`
}

type CreateCustomer struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type UserDB struct {
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type UpdateUser struct {
	Email       string `json:"email" validate:"email"`
	OldPassword string `json:"old_password" binding:"required" validate:"password"`
	NewPassword string `json:"new_password" binding:"required" validate:"password"`
}
type ResponseUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	CreatedAt MyTime `json:"created_at"`
	Role      string `json:"role"`
}

type MockUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email" `
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

// Users array of User type

type Users []User

var PasswordNumber = []rune("0123456789")

var PasswordUpper = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

var PasswordLower = []rune("abcdefghijklmnopqrstuvwxyz")

var PasswordSpecial = []rune("@#%&!$")

var PasswordComposition = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"0123456789" +
	"@#%&!$")

type ErrorResponse struct {
	Message string `json:"message"`
}
