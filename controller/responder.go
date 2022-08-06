package controller

import "github.com/gin-gonic/gin"

type Responder struct{}

func NewResponder() *Responder {
	return &Responder{}
}

func (s Responder) Success(data interface{}) gin.H {
	return gin.H{
		"meta": gin.H{
			"success": true,
		},
		"data": data,
	}
}

func (s Responder) SuccessList(total int, limit int, offset int, data any) gin.H {
	return gin.H{
		"meta": gin.H{
			"success": true,
			"total":   total,
			"limit":   limit,
			"offset":  offset,
		},
		"data": data,
	}
}

func (s Responder) fail(data any) gin.H {
	return gin.H{
		"meta": gin.H{
			"success": false,
		},
		"data": data,
	}
}
