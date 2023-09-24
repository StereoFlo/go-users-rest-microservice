package infrastructure

import "github.com/gin-gonic/gin"

type Responder struct {
	Meta map[string]any `json:"meta"`
	Data any            `json:"data"`
}

func NewResponder() *Responder {
	return &Responder{}
}

func (s Responder) Success(data interface{}) Responder {
	s.Meta = gin.H{
		"success": true,
	}
	s.Data = data

	return s
}

func (s Responder) SuccessList(total int, limit int, offset int, data any) Responder {
	s.Meta = gin.H{
		"success": true,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	}
	s.Data = data

	return s
}

func (s Responder) Fail(data any) Responder {
	s.Meta = gin.H{
		"success": false,
	}
	s.Data = data

	return s
}
