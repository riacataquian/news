package httperror

type Response struct {
	Code int         `json:"statusCode"`
	Data interface{} `json:"data"`
}
