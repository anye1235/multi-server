package mongodb

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"ty/car-prices-master/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ Repo = (*Mongo)(nil)

type Repo interface {
	i()
	Dispose() error
	GetClient() *mongo.Client
	SelectCollection(db, cl string) *mongo.Collection
	FindOneWithOpt(ctx context.Context, db, cl string, filter interface{}, data interface{}, opts ...*options.FindOneOptions) error
	FindOne(ctx context.Context, db, cl string, filter interface{}, data interface{}) error
	Insert(ctx context.Context, db, cl string, data interface{}) (*mongo.InsertOneResult, error)
	FindWithoutPage(ctx context.Context, db, cl string, filter interface{}, sort interface{}, data interface{}) error
	FindManyCommon(ctx context.Context, db, cl string, filter interface{}, data interface{}, options *options.FindOptions) error
}

type Mongo struct {
	client    *mongo.Client
	ctx       context.Context
	defaultDB string
}

func New() (Repo, error) {
	cfg := configs.Get().Mongodb
	address := cfg.Addr
	// db := p.Config.GetMapString("mongo", "db")
	url := fmt.Sprintf("mongodb://%v", address)
	clientOptions := options.Client().ApplyURI(url)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	log.Print("mongo start done: addrss=", clientOptions.Hosts)
	defautlDB := ""
	if nil != cfg.DBs && 0 != len(cfg.DBs) {
		defautlDB = cfg.DBs[0]
	}
	return &Mongo{
		ctx:       ctx,
		client:    client,
		defaultDB: defautlDB,
	}, nil
}

func (p *Mongo) i() {}

func (p *Mongo) Dispose() error {
	if err := p.client.Disconnect(p.ctx); err != nil {
		return err
	}

	return nil
}

func (p *Mongo) GetClient() *mongo.Client {
	return p.client
}

func (p *Mongo) SelectCollection(db, cl string) *mongo.Collection {
	if len(db) == 0 {
		db = p.defaultDB
	}
	return p.client.Database(db).Collection(cl)
}

func (p *Mongo) Insert(ctx context.Context, db, cl string, data interface{}) (*mongo.InsertOneResult, error) {

	if len(db) == 0 {
		db = p.defaultDB
	}
	idValue := reflect.ValueOf(primitive.NewObjectID())
	reflect.ValueOf(data).Elem().FieldByName("ID").Set(idValue)

	result, err := p.client.Database(db).Collection(cl).InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Mongo) InsertMany(ctx context.Context, db, cl string, data []interface{}) (*mongo.InsertManyResult, error) {
	result, err := p.client.Database(db).Collection(cl).InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Mongo) FindOne(ctx context.Context, db, cl string, filter interface{}, data interface{}) error {
	if len(db) == 0 {
		db = p.defaultDB
	}
	return p.client.Database(db).Collection(cl).FindOne(ctx, filter).Decode(data)
}

func (p *Mongo) FindOneWithOpt(ctx context.Context, db, cl string, filter interface{}, data interface{}, opts ...*options.FindOneOptions) error {
	if len(db) == 0 {
		db = p.defaultDB
	}
	return p.client.Database(db).Collection(cl).FindOne(ctx, filter, opts...).Decode(data)
}

func (p *Mongo) FindMany(ctx context.Context, db, cl string, filter interface{}, limit, page int64, sort interface{}, data interface{}) (*Pagination, error) {

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(limit * (page - 1))

	if sort == nil {
		sort = bson.D{{"created_at", -1}}
	}
	options.SetSort(sort)

	curCount, err := p.client.Database(db).Collection(cl).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	count := int64(curCount.RemainingBatchLength())

	pageCount := count / limit
	if count%limit != 0 {
		pageCount++
	}

	cur, err := p.client.Database(db).Collection(cl).Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, data); err != nil {
		return nil, err
	}

	pagination := &Pagination{
		Currentpage: page,
		PageSize:    limit,
		Total:       count,
		TotalPage:   pageCount,
	}

	return pagination, nil
}

func (p *Mongo) FindManyToC(ctx context.Context, db, cl string, filter interface{}, limit, page int64, sort interface{}, data interface{}) (*PaginationToC, error) {

	if limit == 0 {
		limit = 10
	}

	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(limit * page)
	if sort == nil {
		sort = bson.D{{"created_at", -1}}
	}
	options.SetSort(sort)

	curCount, err := p.client.Database(db).Collection(cl).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	count := int64(curCount.RemainingBatchLength())

	pageCount := count / limit
	if count%limit != 0 {
		pageCount++
	}

	cur, err := p.client.Database(db).Collection(cl).Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, data); err != nil {
		return nil, err
	}

	canPull := false
	if page+1 < pageCount {
		canPull = true
	}

	pagination := &PaginationToC{
		Currentpage: page,
		PageSize:    limit,
		CanPull:     canPull,
	}

	return pagination, nil
}

// FindManyCommon
// 查询多个
func (p *Mongo) FindManyCommon(ctx context.Context, db, cl string, filter interface{}, data interface{}, options *options.FindOptions) error {
	if len(db) == 0 {
		db = p.defaultDB
	}
	cur, err := p.client.Database(db).Collection(cl).Find(ctx, filter, options)
	if err != nil {
		return err
	}
	if err := cur.All(ctx, data); err != nil {
		return err
	}
	return nil
}

func (p *Mongo) DeleteOne(ctx context.Context, db, cl string, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := p.client.Database(db).Collection(cl).DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (p *Mongo) DeleteMany(ctx context.Context, db, cl string, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := p.client.Database(db).Collection(cl).DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (p *Mongo) Update(ctx context.Context, db, cl string, filter interface{}, data interface{}, op *options.UpdateOptions) (*mongo.UpdateResult, error) {

	//userInfo, ok := ctx.Value(rest.USERCOOKIE).(*auth.UserInfo)
	//if !ok {
	//	return nil, fmt.Errorf("can not find user")
	//}
	//
	//userNameValue := reflect.ValueOf(userInfo.RealName)
	nowValue := reflect.ValueOf(time.Now().Unix())
	reflect.ValueOf(data).Elem().FieldByName("UpdatedAt").Set(nowValue)
	//reflect.ValueOf(data).Elem().FieldByName("UpdatedBy").Set(userNameValue)

	update := bson.M{"$set": data}
	result, err := p.client.Database(db).Collection(cl).UpdateOne(ctx, filter, update, op)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Mongo) UpdateFiled(ctx context.Context, db, cl string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := p.client.Database(db).Collection(cl).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Mongo) Transaction(ctx context.Context, callback func(mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	// wcMajority := writeconcern.New()

	session, err := p.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Mongo) FindWithoutPage(ctx context.Context, db, cl string, filter interface{}, sort interface{}, data interface{}) error {
	if len(db) == 0 {
		db = p.defaultDB
	}

	options := options.Find()

	options.SetSort(sort)

	cur, err := p.client.Database(db).Collection(cl).Find(ctx, filter, options)
	if err != nil {
		return err
	}

	if err := cur.All(ctx, data); err != nil {
		return err
	}

	return nil
}

func CreateProjectBson(fields ...string) bson.D {
	project := bson.D{}
	for _, field := range fields {
		if len(field) == 0 {
			continue
		}
		project = append(project, bson.E{field, 1})
	}
	return project
}

// func (p *Mongo) Collection(collectionName string) *mongo.Collection {
// 	return p.db.Collection("student")
// }

type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
	DeletedAt time.Time          `json:"DeletedAt" bson:"deleted_at"`
}

type Pagination struct {
	Currentpage int64 `json:"current_page"`
	PageSize    int64 `json:"page_size"`
	Total       int64 `json:"total"`
	TotalPage   int64 `json:"total_page"`
}

type PaginationToC struct {
	Currentpage int64 `json:"current_page"`
	PageSize    int64 `json:"page_size"`
	CanPull     bool
}
