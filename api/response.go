package api

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}
