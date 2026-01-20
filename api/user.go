package api

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name"`  // 用户名
	Email string `json:"email"` // 邮箱
}

// CreateUserResponse 创建用户响应
type CreateUserResponse struct {
	ID int64 `json:"id"` // 用户ID
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name  string `json:"name"`  // 用户名
	Email string `json:"email"` // 邮箱
}

// UserResponse 用户响应
type UserResponse struct {
	ID    int64  `json:"id"`    // 用户ID
	Name  string `json:"name"`  // 用户名
	Email string `json:"email"` // 邮箱
}

// ListUserResponse 用户列表响应
type ListUserResponse struct {
	Total int64           `json:"total"` // 总数
	List  []*UserResponse `json:"list"`  // 用户列表
}

// MessageResponse 消息响应
type MessageResponse struct {
	Message string `json:"message"` // 消息内容
}
