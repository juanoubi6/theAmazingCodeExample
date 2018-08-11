package common

import "github.com/gin-gonic/gin"

type RestError struct {
	Err  error
	Msg  string
	Code int
}

func (r *RestError) ToH() gin.H {
	return gin.H{"msg": r.Msg, "detail": r.Err.Error()}
}
