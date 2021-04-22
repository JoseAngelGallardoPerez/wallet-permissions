package types

type Category struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Sort        int64     `json:"sort"`
	Permissions []*Action `json:"permissions"`
}

type Action struct {
	ActionState
	Id         int64     `json:"id"`
	ParentId   *int64    `json:"parentId"`
	CategoryId *int64    `json:"categoryId"`
	CreatedAt  string    `json:"createdAt"`
	Name       string    `json:"name"`
	Sort       int64     `json:"sort"`
	Children   []*Action `json:"children"`
}
