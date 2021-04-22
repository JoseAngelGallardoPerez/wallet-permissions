package acl

type Privilege string

const (
	PrivilegeRead   = Privilege("read")
	PrivilegeEdit   = Privilege("edit")
	PrivilegeDelete = Privilege("delete")
)

func (p Privilege) String() string {
	return string(p)
}
