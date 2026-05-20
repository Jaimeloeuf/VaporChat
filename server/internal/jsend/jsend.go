package jsend

type Response[T any] struct {
	// enum of "success" | "fail" | "error"
	Status string `json:"status"`

	// This is optional so pointers are used here as nil pointers are omitted
	// during JSON marshalling.
	Data *T `json:"data,omitempty"`

	// This is only used for errors.
	// This is optional so pointers are used here as nil pointers are omitted
	// during JSON marshalling.
	Message *string `json:"message,omitempty"`

	// This is only used for errors.
	// This is optional so pointers are used here as nil pointers are omitted
	// during JSON marshalling.
	Code *int `json:"code,omitempty"`
}

func Success[T any](data T) Response[T] {
	return Response[T]{
		Status: "success",
		Data:   &data,
	}
}

func Fail[T any](data T) Response[T] {
	return Response[T]{
		Status: "fail",
		Data:   &data,
	}
}

func Error(message string) Response[any] {
	return Response[any]{
		Status:  "error",
		Message: &message,
		Code:    nil,
	}
}
