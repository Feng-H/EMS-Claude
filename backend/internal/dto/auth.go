package dto

// Auth DTOs
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token             string  `json:"token"`
	ExpireAt          int64   `json:"expire_at"`
	UserInfo          UserInfo `json:"user_info"`
	MustChangePassword bool   `json:"must_change_password"`
}

type UserInfo struct {
	ID                uint   `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	Role              string `json:"role"`
	FactoryID         *uint  `json:"factory_id"`
	ApprovalStatus    string `json:"approval_status"`
	MustChangePassword bool   `json:"must_change_password"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type RefreshTokenResponse struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
}

// Change Password DTOs
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// User Application DTOs
type UserApplicationRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	Name      string `json:"name" binding:"required,max=100"`
	Role      string `json:"role" binding:"required,oneof=admin supervisor engineer maintenance operator"`
	FactoryID *uint  `json:"factory_id"`
	Phone     string `json:"phone" binding:"omitempty,max=20"`
}

type UserApplicationResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// User Management DTOs
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	Name      string `json:"name" binding:"required,max=100"`
	Role      string `json:"role" binding:"required,oneof=admin supervisor engineer maintenance operator"`
	FactoryID *uint  `json:"factory_id"`
	Phone     string `json:"phone" binding:"omitempty,max=20"`
}

type UpdateUserRequest struct {
	Name      string `json:"name" binding:"omitempty,max=100"`
	Role      string `json:"role" binding:"omitempty,oneof=admin supervisor engineer maintenance operator"`
	FactoryID *uint  `json:"factory_id"`
	Phone     string `json:"phone" binding:"omitempty,max=20"`
	IsActive  *bool  `json:"is_active"`
}

type ApproveUserRequest struct {
	Approve bool   `json:"approve"`
	Reason  string `json:"reason,omitempty"`
}

type UserListResponse struct {
	ID                uint   `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	Role              string `json:"role"`
	Phone             string `json:"phone"`
	IsActive          bool   `json:"is_active"`
	ApprovalStatus    string `json:"approval_status"`
	MustChangePassword bool   `json:"must_change_password"`
	FactoryID         *uint  `json:"factory_id"`
	CreatedAt         string `json:"created_at"`
}
