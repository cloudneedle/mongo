package mongo

import (
	"context"
	"github.com/cloudneedle/mongo/pipe"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"
	"time"
)

var cli *Client
var db *DB

func init() {
	var err error
	cli, err = NewClient(context.Background(), os.Getenv("MONGO"))
	if err != nil {
		panic(err)
	}
	db = cli.NewDB("os_test")
}

type Admin struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
}

func TestAggregate_Pipe(t *testing.T) {
	var md []Admin
	match := pipe.BuildMatcher()
	match.IF(true, bson.E{"_id", ""})
	err := db.Coll("admin").Aggregate(context.Background(), match.Build()).Find(&md)
	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", md)
}

type AppItem struct {
	Id          string  `bson:"_id" json:"id"`    // 应用id
	Name        string  `bson:"name" json:"name"` // 应用名称
	UniqueId    string  `bson:"unique_id" json:"unique_id"`
	DesktopIcon string  `json:"desktop_icon" bson:"desktop_icon"`
	Price       float64 `bson:"price" json:"price" `
	IsBuy       bool    `bson:"is_buy" json:"is_buy"`
}

type PkgItem struct {
	Id        string    `json:"id" bson:"_id"`              // id
	Name      string    `json:"name" bson:"name"`           // 名称
	Price     float64   `json:"price" bson:"price"`         // 价格
	VipPrice  float64   `json:"vip_price" bson:"vip_price"` // 会员价
	Discount  float64   `json:"discount" bson:"discount"`   // 折扣
	AppIds    []string  `json:"app_ids" bson:"app_ids"`     // 应用id
	AppCount  int       `json:"app_count" bson:"app_count"` // 应用数
	Apps      []AppItem `json:"apps" bson:"apps"`           // 应用
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	IsBuy     bool      `bson:"is_buy" json:"is_buy"` // 是否已购买
}

func TestAggregate_FindOne(t *testing.T) {
	lookup := pipe.Lookup(pipe.LookupField{
		From:         "app",
		LocalField:   "app_ids",
		ForeignField: "_id",
		As:           "apps",
	})

	match := pipe.Match(bson.E{Key: "_id", Value: "6389b494721caa0a79bbf41d9"})
	var item PkgItem
	err := db.Coll("app.pkg").Aggregate(context.Background(), lookup, match).FindOne(&item)
	if err != nil {
		t.Error(err)
	}
	t.Log(item)

}
