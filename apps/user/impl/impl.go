package impl

import (
	"context"

	"github.com/zhou-lincong/keyauth/apps/user"
	"github.com/zhou-lincong/keyauth/conf"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &impl{}
)

type impl struct {
	col *mongo.Collection
	log logger.Logger
	user.UnimplementedServiceServer
}

func (s *impl) Config() error {
	// 依赖MongoDB的DB对象
	db, err := conf.C().Mongo.GetDB()
	if err != nil {
		return err
	}

	// 获取一个Collection对象, 通过Collection对象 来进行CRUD
	s.col = db.Collection(s.Name())
	s.log = zap.L().Named(s.Name())

	// 创建索引，避免相同用户注册两次
	indexs := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				// domian和name都在data里面
				{Key: "data.domain", Value: bsonx.Int32(-1)},
				{Key: "data.name", Value: bsonx.Int32(-1)},
			},
			//创建唯一键索引
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = s.col.Indexes().CreateMany(context.Background(), indexs)
	if err != nil {
		return err
	}
	return nil
}

func (s *impl) Name() string {
	return user.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	user.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
