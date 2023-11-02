package model

// A Response to hold meta data and requested data
type Response struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

// A response for if an error occurs
type ErrorResponse struct {
	Status int `json:"status"`
	Name string `json:"name"`
}