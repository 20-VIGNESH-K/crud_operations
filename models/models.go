package models

type Profile struct {
	Name     string `json:"name" bson:"name" validate:"required"`
	Age      int    `json:"age" bson:"age" validate:"gte=0"`
	Address  string `json:"address" bson:"address"`
	Place    string `json:"place" bson:"place"`
	District string `json:"district" bson:"district"`
}
