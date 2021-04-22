package acl

type Resource string

const (
	ResourcePermission = Resource("permission")
	ResourceGroup      = Resource("group")
	ResourceAction     = Resource("action")
)

func (r Resource) String() string {
	return string(r)
}

func (r Resource) ACLType() string {
	return r.String()
}
