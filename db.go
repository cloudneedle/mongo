package mg

import "go.mongodb.org/mongo-driver/mongo"

type DB struct {
	db *mongo.Database
}

func (d *DB) Coll(name string) *Collection {
	return &Collection{coll: d.db.Collection(name)}
}

func (d *DB) Model(coll ICollection) *Collection {
	return &Collection{coll: d.db.Collection(coll.CollName())}
}
