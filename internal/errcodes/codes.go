package errcodes

import "net/http"

const (
	CodeForbidden             = "FORBIDDEN"
	CodeNumeric               = "NUMERIC"
	CodeGroupDoesNotExist     = "GROUP_DOES_NOT_EXIST"
	CodeAdminGroupConstraint  = "ADMIN_GROUP_CONSTRAINT"
	CodeActionKeyDoesNotExist = "ACTION_KEY_DOES_NOT_EXIST"
	CodeUserNotInGroup        = "USER_NOT_IN_GROUP"
	CodeGroupNameDuplication  = "GROUP_NAME_DUPLICATION"
)

var statusCodes = map[string]int{
	CodeForbidden:             http.StatusForbidden,
	CodeNumeric:               http.StatusBadRequest,
	CodeGroupDoesNotExist:     http.StatusBadRequest,
	CodeAdminGroupConstraint:  http.StatusConflict,
	CodeActionKeyDoesNotExist: http.StatusBadRequest,
	CodeUserNotInGroup:        http.StatusBadRequest,
}
