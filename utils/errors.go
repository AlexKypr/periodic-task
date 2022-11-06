package utils

type ErrorDto struct {
	Status string `json:"status"`
	Desc   string `json:"desc"`
}

func NewErrorDTO(status, desc string) ErrorDto {
	return ErrorDto{
		Status: status,
		Desc:   desc,
	}
}
