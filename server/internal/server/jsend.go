package server

type JSendResponse[T any] struct {
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

func JSendSuccess[T any](data T) JSendResponse[T] {
	return JSendResponse[T]{
		Status: "success",
		Data:   &data,
	}
}

func JSendFail[T any](data T) JSendResponse[T] {
	return JSendResponse[T]{
		Status: "fail",
		Data:   &data,
	}
}

func JSendError(message string) JSendResponse[any] {
	return JSendResponse[any]{
		Status:  "error",
		Message: &message,
		Code:    nil,
	}
}
