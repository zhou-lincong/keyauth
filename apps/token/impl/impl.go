package impl

import (
	"context"

	"github.com/zhou-lincong/keyauth/apps/token"
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
	token.UnimplementedServiceServer

	user user.ServiceServer
}

func (s *impl) Config() error {
	// 依赖MongoDB的DB对象
	db, err := conf.C().Mongo.GetDB()
	if err != nil {
		return err
	}

	// 获取一个Collection对象, 通过Collection对象 来进行CRUD
	// 在库里面创建一个token Collection 如果数据库不存在，则创建数据库，否则切换到指定数据库。
	s.col = db.Collection(s.Name()) //集合获取具有给定名称的集合的句柄，该集合的名称配置了给定的 CollectionOptions。
	s.log = zap.L().Named(s.Name())
	// 将依赖引进来
	s.user = app.GetGrpcApp(user.AppName).(user.ServiceServer)

	// 创建索引
	indexs := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "refresh_token", Value: bsonx.Int32(-1)},
			},
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
	return token.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	token.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
