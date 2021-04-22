package acl

type Role int
type Roles []Role
type rolesHelper struct{}

var (
	roles = map[int]string{
		-1:         "undefined",
		128:        "guest",
		32768:      "client",
		8388608:    "admin",
		2147483648: "root",
	}
	indexes = map[string]int{
		"undefined": -1,
		"guest":     128,
		"client":    32768,
		"admin":     8388608,
		"root":      2147483648,
	}
)

var RolesHelper = rolesHelper{}

const (
	Undefined Role = -1
	Guest     Role = 128
	Client    Role = 32768
	Buyer     Role = 32769
	Supplier  Role = 32770
	Financier Role = 32771
	Admin     Role = 8388608
	Root      Role = 2147483648
)

func (r Role) String() string {
	return roles[int(r)]
}

func (r Role) ACLType() string {
	return r.String()
}

func (r Role) Is(role Role) bool {
	return RolesHelper.FromName(role.String()) == RolesHelper.FromName(r.String())
}

func (*rolesHelper) FromName(role string) Role {
	if idx, ok := indexes[role]; ok {
		return Role(idx)
	}
	return Role(-1)
}

func (r Roles) Strings() []string {
	result := make([]string, len(r))
	for i, role := range r {
		result[i] = role.String()
	}
	return result
}
