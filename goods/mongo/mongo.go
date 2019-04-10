package mongo

import (
	"context"
	"errors"
	//"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"net/url"
	"reflect"
	"rubik/server/common/t"
	"time"
)

type mongodb struct {
	client *mongo.Client
	db     *mongo.Database
	dbname string
}

var Client *mongodb
var collection *mongo.Collection

func Init(dbUrl string) {
	Client = New(dbUrl)
}

func New(urlStr string) *mongodb {
	dbUrl, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	var mgo mongodb
	//fmt.Println(dbUrl)
	mgo.dbname = dbUrl.Path[1:]
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, urlStr)
	if err != nil {
		panic(err)
	}
	mgo.db = client.Database(mgo.dbname)
	mgo.client = client
	return &mgo
}


func (this *mongodb) GoodsFind(table string,q interface{} , result interface{}) (err error) {
	cursor, err := this.client.Database("goods").Collection(table).Find(context.Background(), q)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		//data := bson.D{}
		var data interface{}
		if err = cursor.Decode(&data); err != nil {
			return err
		}
		*result.(*[]map[string]interface{}) = append(*result.(*[]map[string]interface{}), this.Map(data))
	}
	return
}

func (this *mongodb) Count(table string, q interface{}) (int64, error) {
	count, err := this.client.Database("goods").Collection(table).CountDocuments(context.Background(), q)
	return count, err
}


func (this *mongodb) Insert(table string, data interface{}) error {
	_, err := this.client.Database("goods").Collection(table).InsertOne(context.Background(), data)
	return err
}

func (this *mongodb) Aggregate(table string, pipeline interface{}, result interface{}) (err error) {
	cursor, err := this.client.Database("goods").Collection(table).Aggregate(context.Background(), pipeline)
	if err != nil {
		return
	}
	if err = cursor.Err(); err != nil {
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
	var doc interface{}
	if err = cursor.Decode(&doc); err != nil {
		return
	}
	*result.(*[]map[string]interface{}) = append(*result.(*[]map[string]interface{}), this.Map(doc))
	}
		return
}

func (this *mongodb) Update(table string, q interface{}, update interface{}) error {
	_, err := this.client.Database("goods").Collection(table).UpdateOne(context.Background(), q, update)
	return err
}


















func (this *mongodb) Array(result interface{}) []interface{} {
	resultArray := result.(primitive.A)
	for k, v := range resultArray {
		if v == nil {
			continue
		}
		switch reflect.TypeOf(v).String() {
		case "primitive.D":
			resultArray[k] = this.Map(v)
		case "primitive.ObjectID":
			resultArray[k] = v.(primitive.ObjectID).Hex()
		case "primitive.A":
			resultArray[k] = this.Array(v)
		}
	}
	return resultArray
}
func (this *mongodb) Map(result interface{}) map[string]interface{} {
	resultMap := result.(primitive.D).Map()
	for k, v := range resultMap {
		if v == nil {
			continue
		}
		switch reflect.TypeOf(v).String() {
		case "primitive.D":
			resultMap[k] = this.Map(v)
		case "primitive.ObjectID":
			resultMap[k] = v.(primitive.ObjectID).Hex()
		case "primitive.A":
			resultMap[k] = this.Array(v)
		}
	}
	return resultMap
}
func (this *mongodb) FindOne(table string, q interface{}, result interface{}, orther ...interface{}) (err error) {
	var opts []*options.FindOneOptions
	for _, opt := range orther {
		one := options.FindOne()
		for k, v := range opt.(bson.M) {
			if k == "Sort" {
				one.Sort = v
			}
		}
		opts = append(opts, one)
	}

	var data primitive.D
	if err = this.client.Database(this.dbname).Collection(table).FindOne(context.Background(), q, opts...).Decode(&data); err != nil {
		return
	}
	//多3倍性能
	*result.(*map[string]interface{}) = this.Map(data)
	return
}

func (this *mongodb) Find(table string, q interface{}, result interface{}, orther ...interface{}) (err error) {
	var opts []*options.FindOptions
	var collOpts []*options.CollectionOptions
	for _, opt := range orther {
		one := options.Find()
		for k, v := range opt.(bson.M) {
			if k == "Sort" {
				one.Sort = v
			} else if k == "Limit" {
				i := t.Int64(v)
				one.Limit = &i
			} else if k == "Skip" {
				i := t.Int64(v)
				one.Skip = &i
			} else if k == "ReadPreference" {
				var ref *readpref.ReadPref
				if v == "SecondaryPreferred" {
					ref = readpref.SecondaryPreferred()
				} else if v == "Secondary" {
					ref = readpref.Secondary()
				}
				collOpts = []*options.CollectionOptions{&options.CollectionOptions{ReadPreference: ref}}
			}
		}
		opts = append(opts, one)
	}

	cursor, err := this.client.Database(this.dbname).Collection(table, collOpts...).Find(context.Background(), q, opts...)
	if err != nil {
		return
	}
	if err = cursor.Err(); err != nil {
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		//data := bson.D{}
		var data interface{}
		if err = cursor.Decode(&data); err != nil {
			return
		}
		*result.(*[]map[string]interface{}) = append(*result.(*[]map[string]interface{}), this.Map(data))
	}
	return
}
func (this *mongodb) Index(table string, models map[string]interface{}) error {
	var unique bool
	if models["unique"] == nil {
		unique = false
	} else {
		unique = models["unique"].(bool)
	}
	_, err := this.client.Database(this.dbname).Collection(table).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: models["keys"],
		Options: &options.IndexOptions{
			Unique: &unique,
		}})
	return err
}
func (this *mongodb) Indexs(table string, models []map[string]interface{}) error {
	var modelList = make([]mongo.IndexModel, 0)
	for _, v := range models {
		var unique bool
		if v["unique"] == nil {
			unique = false
		} else {
			unique = v["unique"].(bool)
		}
		modelList = append(modelList, mongo.IndexModel{
			Keys: v["keys"],
			Options: &options.IndexOptions{
				Unique: &unique,
			},
		})
	}
	_, err := this.client.Database(this.dbname).Collection(table).Indexes().CreateMany(context.Background(), modelList)
	return err
}

func (this *mongodb) Update2(session *mongo.SessionContext, table string, q interface{}, update interface{}) error {
	_, err := this.client.Database(this.dbname).Collection(table).UpdateOne(context.Background(), q, update)
	return err
}

func (this *mongodb) Insert2(session *mongo.SessionContext, table string, data interface{}) error {
	_, err := this.client.Database(this.dbname).Collection(table).InsertOne(*session, data)
	return err
}
func (this *mongodb) FindOneAndUpdate(table string, q interface{}, update interface{}, result interface{}, orther ...interface{}) (err error) {
	var opts []*options.FindOneAndUpdateOptions
	for _, opt := range orther {
		one := options.FindOneAndUpdate()
		for k, v := range opt.(bson.M) {
			if k == "Sort" {
				one.Sort = v
			} else if k == "Upsert" {
				i := v.(bool)
				one.Upsert = &i
			} else if k == "New" {
				i := v.(bool)
				if i == true {
					one.SetReturnDocument(options.After)
				} else if i == false {
					one.SetReturnDocument(options.Before)
				}
			}
		}
		opts = append(opts, one)
	}
	var data interface{}
	if err = this.client.Database(this.dbname).Collection(table).FindOneAndUpdate(context.Background(), q, update, opts...).Decode(&data); err != nil {
		return
	}
	*result.(*map[string]interface{}) = this.Map(data)
	return
}
func (this *mongodb) FindOneAndUpdate2(table string, session *mongo.SessionContext, q interface{}, update interface{}, result interface{}, orther ...interface{}) error {
	var opts []*options.FindOneAndUpdateOptions
	for _, opt := range orther {
		one := options.FindOneAndUpdate()
		for k, v := range opt.(bson.M) {
			if k == "Sort" {
				one.Sort = v
			} else if k == "Upsert" {
				i := v.(bool)
				one.Upsert = &i
			} else if k == "New" {
				i := v.(bool)
				if i == true {
					one.SetReturnDocument(options.After)
				} else if i == false {
					one.SetReturnDocument(options.Before)
				}
			}
		}
		opts = append(opts, one)
	}
	var data primitive.D
	err := this.db.Collection(table).FindOneAndUpdate(*session, q, update, opts...).Decode(&data)
	*result.(*map[string]interface{}) = this.Map(data)
	return err
}
func (this *mongodb) GetID(name string) (int64, error) {
	var doc map[string]interface{}
	err := this.FindOneAndUpdate("ids", bson.M{"name": name}, bson.M{"$inc": bson.M{"id": int64(1)}}, &doc, bson.M{"Upsert": true, "New": true})
	return doc["id"].(int64), err
}
func (this *mongodb) DeleteOne(table string, q interface{}) error {
	result, err := this.client.Database(this.dbname).Collection(table).DeleteOne(context.Background(), q)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Not Fond")
	}
	return nil
}
