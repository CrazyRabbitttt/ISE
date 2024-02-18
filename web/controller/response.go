package controller

import "net/http"

type ResponseData struct {
	State   int32       `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseOkWithMessage(data interface{}) *ResponseData {
	return &ResponseData{
		State:   http.StatusOK,
		Message: "success",
		Data:    data,
	}
}
