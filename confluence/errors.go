package confluence

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ErrorResponse 映射 Confluence API 返回的错误结构
type ErrorResponse struct {
	StatusCode int `json:"-"` // HTTP 状态码
	Data       struct {
		Authorized bool     `json:"authorized"`
		Valid      bool     `json:"valid"`
		Errors     []APIError `json:"errors"`
	} `json:"data"`
	Message string `json:"message"` // 顶层错误信息
}

// APIError 具体的错误详情
type APIError struct {
	Message struct {
		Translation string   `json:"translation"`
		Args        []string `json:"args"`
	} `json:"message"`
}

func (e *ErrorResponse) Error() string {
	if len(e.Data.Errors) > 0 {
		return fmt.Sprintf("confluence api error: %d - %s (details: %s)", e.StatusCode, e.Message, e.Data.Errors[0].Message.Translation)
	}
	return fmt.Sprintf("confluence api error: %d - %s", e.StatusCode, e.Message)
}

// CheckResponse 检查 API 响应状态码，如果非 2xx 则返回 error
func CheckResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	errorResponse := &ErrorResponse{StatusCode: res.StatusCode}
	data, err := io.ReadAll(res.Body)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}
