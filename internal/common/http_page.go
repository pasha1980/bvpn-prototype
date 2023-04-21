package common

type HttpPage struct {
	Code  int     `json:"code"`
	Error *string `json:"error"`
	Data  any     `json:"data"`
}

func MakeHttpPage(data any) HttpPage {
	return HttpPage{
		Code:  0,
		Error: nil,
		Data:  data,
	}
}

func MakeErrorPage(code int, error string, data any) HttpPage {
	return HttpPage{
		Code:  code,
		Error: &error,
		Data:  data,
	}
}
