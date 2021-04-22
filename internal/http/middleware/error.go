package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Confialink/wallet-permissions/internal/http/response"
	"github.com/Confialink/wallet-permissions/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	ErrorInternalError = errors.New("something went wrong")
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypePublic:
					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						res := response.New().AddError(
							e.Error(),
							"public",
							nil,
							nil,
							e.Meta,
						)
						c.JSON(c.Writer.Status(), res)
					}
				case gin.ErrorTypeBind:
					res := response.New()
					if validationErr, ok := e.Err.(validator.ValidationErrors); ok {
						for _, err := range validationErr {
							res.AddError(
								ValidationErrorToText(err),
								"validation",
								nil,
								nil,
								gin.H{"field": err.Field},
							)
						}
					}
					// Make sure we maintain the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}
					c.JSON(status, res)

				}
			}

			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() && c.Writer.Status() == http.StatusOK {
				c.JSON(
					http.StatusInternalServerError,
					response.NewWithError(ErrorInternalError.Error(), "internal"),
				)
			}
		}
	}
}

func ValidationErrorToText(e validator.FieldError) string {
	field := util.Text.ToDelimiter(e.Field(), " ")
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required and cannot be empty", field)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", field, e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", field, e.Param())
	case "email":
		return fmt.Sprintf("'%s' is not a valid email address", e.Value())
	case "len":
		return fmt.Sprintf("%s must be %s characters long", field, e.Param())
	case "eqfield":
		return fmt.Sprintf("%s must match the %s", field, e.Param())
	case "pwdpolicy":
		return fmt.Sprintf("%s must contain upper and lower case characters, and at least one digit and special character", e.Field())
	}
	return fmt.Sprintf("%s is not valid", field)
}
