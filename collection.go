package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICollection interface {
	CollName() string
}

type Collection struct {
	coll *mongo.Collection
}

type Filter struct {
}

type Aggregate struct {
	pipes []bson.D
	coll  *mongo.Collection
	ctx   context.Context
}

func (a *Aggregate) Find(data any) error {
	cur, err := a.coll.Aggregate(a.ctx, a.pipes)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return err
	}
	defer cur.Close(a.ctx)
	err = cur.All(a.ctx, data)
	return err
}

func (a *Aggregate) FindOne(data any) error {
	cur, err := a.coll.Aggregate(a.ctx, a.pipes)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return err
	}
	defer cur.Close(a.ctx)
	if cur.TryNext(a.ctx) {
		err = cur.Decode(data)
	}
	return err
}

func (a *Aggregate) FindPageList(pageIndex, pageSize int64, data any) (int64, error) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	count, err := a.Count()
	if err != nil {
		return 0, err
	}

	var pipes = a.pipes
	if count > 0 {
		a.pipes = append(a.pipes, bson.D{{"$skip", (pageIndex - 1) * pageSize}})
		a.pipes = append(a.pipes, bson.D{{"$limit", pageSize}})
	}

	err = a.Find(data)
	if err != nil {
		return 0, err
	}
	a.pipes = pipes
	return count, nil
}

func (a *Aggregate) Count() (int64, error) {
	var pipes = a.pipes
	a.pipes = append(a.pipes, bson.D{{"$count", "count"}})
	var data []map[string]int64
	err := a.Find(&data)
	if err != nil {
		return 0, err
	}
	a.pipes = pipes
	if data != nil && len(data) > 0 {
		return data[0]["count"], nil
	}
	return 0, nil
}

func (c *Collection) FindByID(ctx context.Context, id string, data any) error {
	err := c.coll.FindOne(ctx, bson.M{"_id": id}).Decode(data)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (c *Collection) FindByObjID(ctx context.Context, id primitive.ObjectID, data any) error {
	err := c.coll.FindOne(ctx, bson.M{"_id": id}).Decode(data)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func (c *Collection) Find(ctx context.Context, filter bson.M, data any) error {
	cur, err := c.coll.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, data)
	return err
}

func (c *Collection) FindOne(ctx context.Context, filter bson.M, data any) error {
	err := c.coll.FindOne(ctx, filter).Decode(data)
	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func (c *Collection) InsertOne(ctx context.Context, doc any) error {
	_, err := c.coll.InsertOne(ctx, doc)
	return err
}

func (c *Collection) InsertMany(ctx context.Context, docs []any) error {
	_, err := c.coll.InsertMany(ctx, docs)

	return err
}

func (c *Collection) DeleteOne(ctx context.Context, filter bson.M) error {
	_, err := c.coll.DeleteOne(ctx, filter)
	return err
}

func (c *Collection) DeleteMany(ctx context.Context, filter bson.M) error {
	_, err := c.coll.DeleteMany(ctx, filter)
	return err
}

func (c *Collection) UpdateOne(ctx context.Context, filter bson.M, data any) error {
	_, err := c.coll.UpdateOne(ctx, filter, data)
	return err
}

func (c *Collection) UpdateMany(ctx context.Context, filter bson.M, data any) error {
	_, err := c.coll.UpdateMany(ctx, filter, data)
	return err
}

func (c *Collection) ReplaceOne(ctx context.Context, filter bson.M, data any) error {
	_, err := c.coll.ReplaceOne(ctx, filter, data)
	return err
}

func (c *Collection) Aggregate(ctx context.Context, pipes ...bson.D) *Aggregate {
	var p []bson.D
	if len(pipes) > 0 {
		for _, v := range pipes {
			if len(v) > 0 {
				p = append(p, v)
			}
		}
	}
	return &Aggregate{
		pipes: p,
		coll:  c.coll,
		ctx:   ctx,
	}
}
