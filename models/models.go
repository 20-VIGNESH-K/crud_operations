package models

type Profile struct {
	Name     string `json:"name" bson:"name"`
	Age      int    `json:"age" bson:"age"`
	Address  string `json:"address" bson:"address"`
	Place    string `json:"place" bson:"place"`
	District string `json:"district" bson:"district"`
}
