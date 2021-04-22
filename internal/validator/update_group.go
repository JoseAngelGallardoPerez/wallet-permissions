package validator

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-permissions/internal/db/model"
	"github.com/Confialink/wallet-permissions/internal/errcodes"
	"github.com/Confialink/wallet-permissions/internal/http/form"
	"github.com/Confialink/wallet-permissions/internal/service"
)

func UpdateGroup(c *gin.Context, model *model.Group, groupsService *service.Groups) (*form.UpdateGroup, errors.TypedError) {
	newForm := form.UpdateGroup{}
	if err := c.ShouldBind(&newForm); err != nil {
		return &newForm, errors.ShouldBindToTyped(err)
	}

	if newForm.Name != model.Name() {
		if typedErr := checkGroupNameDuplication(newForm.Name, groupsService); typedErr != nil {
			return &newForm, typedErr
		}
	}

	return &newForm, nil
}

func checkGroupNameDuplication(name string, groupsService *service.Groups) errors.TypedError {
	foundGroup, typedErr := groupsService.FindByName(name)
	if typedErr != nil {
		return typedErr
	}
	if foundGroup.IsExist() {
		return &errors.ValidationErrors{Errors: []errors.ValidationError{
			{
				Code:   errcodes.CodeGroupNameDuplication,
				Source: "name",
				Meta: struct {
					Value string `json:"value"`
				}{name},
			},
		}}
	}
	return nil
}
