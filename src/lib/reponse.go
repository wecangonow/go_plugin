package lib

type Response struct {
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Code   string      `json:"code"`
}


type JsAdRes struct {
	Url     string `json:"url"`
	Ad_id   string `json:"ad_id"`
	Ad_name string `json:"ad_name"`
	Replace string `json:"replace"`
	Ecpm    string `json:"ecpm"`
	Ad_type string `json:"ad_type"`
	Ad_sort string `json:"ad_sort"`
}
