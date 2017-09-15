package modules

type Tbl_user struct {
	ID   string
	NAME string
	AGE int
}

func (t Tbl_user) TableName() string {
	return "gft.tbl_user"
}
