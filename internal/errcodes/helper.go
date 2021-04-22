package errcodes

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
)

func AddError(c *gin.Context, code string) {
	errors.AddErrors(c, CreatePublicError(code))
}

func CreatePublicError(code string) *errors.PublicError {
	return &errors.PublicError{
		Code:       code,
		HttpStatus: statusCodes[code],
	}
}

func AddErrorMeta(c *gin.Context, code string, meta interface{}) {
	publicErr := &errors.PublicError{
		Code:       code,
		HttpStatus: statusCodes[code],
		Meta:       meta,
	}
	errors.AddErrors(c, publicErr)
}
