package app

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-permissions/internal/errcodes"
	"github.com/Confialink/wallet-permissions/internal/service"
)

type innerValidator struct {
	groupsService *service.Groups
	validate      *validator.Validate
	logger        log15.Logger
}

// LoadValidator fetches Validate and registers validations
func LoadValidator(groupsService *service.Groups, logger log15.Logger) *validator.Validate {
	inner := newInnerValidator(groupsService, logger)
	inner.init()
	return inner.validator()
}

func newInnerValidator(groupsService *service.Groups, logger log15.Logger) *innerValidator {
	return &innerValidator{
		groupsService: groupsService,
		validate:      binding.Validator.Engine().(*validator.Validate),
		logger:        logger,
	}
}

func (v *innerValidator) validator() *validator.Validate {
	return v.validate
}

func (v *innerValidator) init() {
	if err := v.validate.RegisterValidation("groupNameUniqueness", v.groupNameUniqueness); err != nil {
		v.logger.Error("Can't register validation", "error", err)
		panic(err)
	}

	errors.SetFormatters(map[string]*errors.ValidationErrorFormatter{"groupNameUniqueness": {
		Code: errcodes.CodeGroupNameDuplication,
		MetaFunc: func(e validator.FieldError) interface{} {
			return struct {
				Value string `json:"value"`
			}{e.Value().(string)}
		},
	}})
}

func (v *innerValidator) groupNameUniqueness(fl validator.FieldLevel) bool {
	name := fl.Field().Interface().(string)
	group, err := v.groupsService.FindByName(name)
	if err != nil {
		v.logger.Error("Can't fetch group", "error", err)
		return false
	}
	return !group.IsExist()
}
