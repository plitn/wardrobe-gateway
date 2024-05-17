package model

type User struct {
	Id       int64  `db:"id" json:"id" goku:"skipinsert"`
	Name     string `db:"name" json:"name"`
	Password string `db:"password" json:"password"`
	ImageUrl string `db:"image_url" json:"image_url"`
	Gender   int    `db:"gender" json:"gender"`
	Age      int    `db:"age" json:"age"`
	Weight   int    `db:"weight" json:"weight"`
}
