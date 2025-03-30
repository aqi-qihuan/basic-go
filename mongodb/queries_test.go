package mongodb

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

type MongoDBTestSuite struct {
	suite.Suite
	col *mongo.Collection
}

// SetupSuite 是 MongoDBTestSuite 的一个方法，用于在测试套件开始前进行初始化设置
func (s *MongoDBTestSuite) SetupSuite() {
	// 获取测试实例
	t := s.T()
	// 创建一个带有超时的上下文，超时时间为10秒
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 在函数结束时取消上下文，释放资源
	defer cancel()

	// 创建一个命令监控器，用于监控 MongoDB 命令的执行
	monitor := &event.CommandMonitor{
		// 当命令开始执行时调用，打印命令内容
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}
	// 创建 MongoDB 客户端选项，设置连接 URI 和命令监控器
	opts := options.Client().
		ApplyURI("mongodb://root:example@localhost:27017/").
		SetMonitor(monitor)
	// 连接到 MongoDB，返回客户端实例和错误信息
	client, err := mongo.Connect(ctx, opts)
	// 断言连接过程中没有错误
	assert.NoError(t, err)
	// 操作 client
	col := client.Database("lmbook").
		Collection("articles")
	s.col = col

	manyRes, err := col.InsertMany(ctx, []any{Article{
		Id:       123,
		AuthorId: 11,
	}, Article{
		Id:       234,
		AuthorId: 12,
	}})

	assert.NoError(s.T(), err)
	s.T().Log("插入数量", len(manyRes.InsertedIDs))
}

// TearDownSuite 是 MongoDBTestSuite 类型的测试套件拆卸方法
// 该方法在测试套件执行完毕后被调用，用于清理测试过程中创建的数据和索引
func (s *MongoDBTestSuite) TearDownSuite() {
	// 创建一个带有超时的上下文，超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 在函数结束时调用cancel函数，释放资源
	defer cancel()
	// 调用DeleteMany方法删除集合中的所有文档
	// bson.D{} 表示空的过滤条件，即匹配所有文档
	_, err := s.col.DeleteMany(ctx, bson.D{})
	// 使用assert.NoError断言错误为nil，如果不是则测试失败
	assert.NoError(s.T(), err)
	// 调用Indexes().DropAll方法删除集合中的所有索引
	_, err = s.col.Indexes().DropAll(ctx)
	// 使用assert.NoError断言错误为nil，如果不是则测试失败
	assert.NoError(s.T(), err)
}

// MongoDBTestSuite 是一个测试套件，用于测试 MongoDB 相关的功能
func (s *MongoDBTestSuite) TestOr() {
	// 创建一个带有超时的上下文，超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 在函数结束时调用cancel函数，释放资源
	defer cancel()
	// 定义一个过滤条件，使用bson数组表示，其中包含两个bson文档
	// 每个bson文档表示一个条件，这里是id等于123或id等于234
	filter := bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"id", 234}}}
	// 调用MongoDB集合的Find方法，传入上下文和过滤条件
	// 过滤条件使用$or操作符，表示查询id为123或234的文档
	res, err := s.col.Find(ctx, bson.D{bson.E{"$or", filter}})
	// 断言错误为nil，如果err不为nil，则测试失败
	assert.NoError(s.T(), err)
	// 定义一个Article类型的切片，用于存储查询结果
	var arts []Article
	// 调用Find方法的All方法，将查询结果解码到arts切片中
	err = res.All(ctx, &arts)
	// 断言错误为nil，如果err不为nil，则测试失败
	assert.NoError(s.T(), err)
	// 使用测试框架的Log方法输出查询结果
	s.T().Log("查询结果", arts)
}

// MongoDBTestSuite 是一个测试套件，用于测试 MongoDB 相关的功能
func (s *MongoDBTestSuite) TestAnd() {
	// 创建一个带有超时的上下文，超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 在函数结束时调用 cancel 函数，释放资源
	defer cancel()
	// 定义一个过滤器，包含两个条件：id 为 123 和 author_id 为 11
	filter := bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"author_id", 11}}}
	// 在集合 s.col 中执行查询，查询条件为 $and 操作符，结合上述两个条件
	res, err := s.col.Find(ctx, bson.D{bson.E{"$and", filter}})
	// 断言查询过程中没有错误发生
	assert.NoError(s.T(), err)
	// 定义一个 Article 类型的切片，用于存储查询结果
	var arts []Article
	// 将查询结果转换为 Article 类型的切片
	err = res.All(ctx, &arts)
	// 断言转换过程中没有错误发生
	assert.NoError(s.T(), err)
	// 打印查询结果
	s.T().Log("查询结果", arts)
}

func (s *MongoDBTestSuite) TestIn() {
	// 创建一个带有超时的上下文，超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 在函数结束时取消上下文，释放资源
	defer cancel()
	// 构建查询过滤器，使用bson.D表示文档，bson.E表示键值对
	// 这里查询id字段在给定数组中的文档
	filter := bson.D{bson.E{"id",
		bson.D{bson.E{"$in", []int{123, 234}}}}}
	// 构建投影，指定只返回id字段，1表示包含该字段
	// 注释掉的代码是另一种表示方式，使用bson.D
	//proj := bson.D{bson.E{"id", 1}}
	proj := bson.M{"id": 1}
	// 执行查询，传入上下文、过滤器和投影选项
	res, err := s.col.Find(ctx, filter,
		options.Find().SetProjection(proj))
	// 断言查询过程中没有错误发生
	assert.NoError(s.T(), err)
	// 定义一个Article类型的切片用于存储查询结果
	var arts []Article
	// 将查询结果解码到arts切片中
	err = res.All(ctx, &arts)
	// 断言解码过程中没有错误发生
	assert.NoError(s.T(), err)
	// 打印查询结果
	s.T().Log("查询结果", arts)
}

// MongoDBTestSuite 是一个测试套件，用于测试 MongoDB 相关的功能
func (s *MongoDBTestSuite) TestIndexes() {
	// 创建一个带有超时的上下文，超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 在函数结束时调用cancel函数，确保上下文被正确取消，释放资源
	defer cancel()
	// 在MongoDB集合中创建一个索引
	// ires 用于存储索引创建的结果
	// err 用于存储可能发生的错误
	ires, err := s.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		// 指定索引的键，这里是一个升序的 "id" 字段
		Keys: bson.D{bson.E{"id", 1}},
		// 设置索引的选项
		// SetUnique(true) 表示索引是唯一的，不允许重复的 "id" 值
		// SetName("idx_id") 设置索引的名称为 "idx_id"
		Options: options.Index().SetUnique(true).SetName("idx_id"),
	})
	// 使用 assert.NoError 检查 err 是否为 nil，如果不是则测试失败
	assert.NoError(s.T(), err)
	// 使用 s.T().Log 输出创建索引的结果
	s.T().Log("创建索引", ires)
}

func TestMongoDBQueries(t *testing.T) {
	suite.Run(t, &MongoDBTestSuite{})
}
