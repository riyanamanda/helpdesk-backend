package response

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type SuccessResponse[T any] struct {
	Data T     `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Success[T any](data T, meta *Meta) SuccessResponse[T] {
	return SuccessResponse[T]{
		Data: data,
		Meta: meta,
	}
}

func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}
