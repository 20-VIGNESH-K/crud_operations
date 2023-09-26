package models

type Profile struct {
	Name     string `json:"name" bson:"name" validate:"required,customValidator"`
	Age      int    `json:"age" bson:"age" validate:"required,gte=1,lte=100"`
	Address  string `json:"address" bson:"address" validate:"required"`
	Place    string `json:"place" bson:"place" validate:"required"`
	District string `json:"district" bson:"district" validate:"required"`
}
