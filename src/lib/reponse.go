package lib

type Response struct {
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Code   string      `json:"code"`
}
