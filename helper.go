package mg

import "go.mongodb.org/mongo-driver/bson/primitive"

func GetStrId() string {
	return primitive.NewObjectID().Hex()
}
