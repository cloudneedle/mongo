package pipe

import "go.mongodb.org/mongo-driver/bson"

func Match(items ...bson.E) bson.D {
	if len(items) > 0 {
		return bson.D{{"$match", items}}
	}
	return nil
}

type Matcher struct {
	container bson.D
}

func BuildMatcher() *Matcher {
	return &Matcher{}
}

func (m *Matcher) IF(condition bool, e bson.E) *Matcher {
	if condition {
		if len(e.Key) > 0 {
			m.container = append(m.container, e)
		}
	}
	return m
}

func (m *Matcher) Match(e bson.E) *Matcher {
	if len(e.Key) > 0 {
		m.container = append(m.container, e)
	}
	return m
}

func (m *Matcher) Build() bson.D {
	if len(m.container) > 0 {
		return bson.D{{"$match", m.container}}
	}
	return nil
}

type LookupField struct {
	From         string `bson:"from"`
	LocalField   string `bson:"localField"`
	ForeignField string `bson:"foreignField"`
	As           string `bson:"as"`
}

// Lookup 查询
func Lookup(field LookupField) bson.D {
	return bson.D{
		{"$lookup",
			bson.D{{"from", field.From},
				{"localField", field.LocalField},
				{"foreignField", field.ForeignField},
				{"as", field.As}}},
	}
}

type UnwindField struct {
	Path                       string `json:"path"`
	IncludeArrayIndex          string `json:"includeArrayIndex"`
	PreserveNullAndEmptyArrays bool   `json:"preserveNullAndEmptyArrays"`
}

// Unwind 查询
func Unwind(field UnwindField) bson.D {
	val := bson.D{
		{"path", field.Path},
	}
	if field.IncludeArrayIndex != "" {
		val = append(val, bson.E{Key: "includeArrayIndex", Value: field.IncludeArrayIndex})
	}
	if field.PreserveNullAndEmptyArrays {
		val = append(val, bson.E{Key: "preserveNullAndEmptyArrays", Value: field.PreserveNullAndEmptyArrays})
	}
	return bson.D{
		{"$unwind", val},
	}
}

// Group 分组
func Group(items ...bson.E) bson.D {
	if len(items) > 0 {
		return bson.D{{"$group", items}}
	}
	return nil
}

// Project 选择
func Project(items ...bson.E) bson.D {
	if len(items) > 0 {
		return bson.D{{"$project", items}}
	}
	return nil
}

// Sort 排序
func Sort(items ...bson.E) bson.D {
	if len(items) > 0 {
		return bson.D{{"$sort", items}}
	}
	return nil
}
