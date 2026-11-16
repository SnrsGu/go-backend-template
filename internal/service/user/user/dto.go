package user

type CreateUserRequest struct {
	Username string `json:"username"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type UserListResponse struct {
	Items    []UserResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}
