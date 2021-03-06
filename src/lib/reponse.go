package lib

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Code string      `json:"code"`
}

type JsAdRes struct {
	Url     string  `json:"url"`
	Ad_id   int     `json:"ad_id"`
	Ad_name string  `json:"ad_name"`
	Replace int     `json:"replace"`
	Ecpm    float64 `json:"ecpm"`
	Ad_type int     `json:"ad_type"`
	Ad_sort string  `json:"AD_sort"`
}

type DeployAdRes struct {
	Url     string  `json:"url"`
	Ad_id   int     `json:"ad_id"`
	Ad_name string  `json:"ad_name"`
	Web     string  `json:"web"`
	Ecpm    float64 `json:"ecpm"`
	Ad_type int     `json:"ad_type"`
	Ad_sort string  `json:"AD_sort"`
	Size    string  `json:"size"`
}

type FinalResponse struct {
	JsAd     JsAdRes       `json:"jsAd"`
	DeployAd []DeployAdRes `json:"deployAd"`
	UrlType  int           `json:"URL_type"`
}
