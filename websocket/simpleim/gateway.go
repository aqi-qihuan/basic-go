package simpleim

import (
	"basic-go/lmbook/pkg/logger"
	"basic-go/lmbook/pkg/saramax"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
)

// websocket 的网关
type WsGateway struct {
	// 连接了这个实例的客户端
	// 这里我们用 uid 作为 key
	// 实践中要考虑到不同的设备，
	// 那么这个 key 可能是一个复合结构，例如 uid + 设备
	conns *syncx.Map[int64, *Conn]
	svc   *IMService

	client     sarama.Client
	instanceId string
}

func (g *WsGateway) Start(addr string) error {
	// 接收 websocket
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", g.wsHandler)
	// 监听别的节点转发的消息
	err := g.subscribeMsg()
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, mux)
}

func (g *WsGateway) subscribeMsg() error {
	cg, err := sarama.NewConsumerGroupFromClient(g.instanceId,
		g.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{eventName},
			saramax.NewHandler[Event](logger.NewNoOpLogger(), g.consume))
		if err != nil {
			log.Println("退出监听消息循环", err)
		}
	}()
	return nil
}

func (g *WsGateway) wsHandler(writer http.ResponseWriter, request *http.Request) {
	upgrader := websocket.Upgrader{}
	uid := g.Uid(request)
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		writer.Write([]byte("初始化 websocket 失败"))
		return
	}
	c := &Conn{Conn: conn}
	g.conns.Store(uid, c)
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("接收 websocket 数据失败", err)
				return
			}
			// 转发到后端
			var msg Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("非法数据格式", err)
				continue
			}
			// 通知后端
			go func() {
				// 这里要开 goroutine，因为一条消息的处理过程，可能很慢
				ctx, cancel := context.WithTimeout(context.Background(),
					time.Second*3)
				defer cancel()
				err1 := g.svc.Receive(ctx, uid, msg)
				if err1 != nil {
					// 正常来说，这里出错的时候，要通知用户发送失败
					err1 = c.Send(Message{Seq: msg.Seq, Type: "result", Content: "FAILED"})
					// 这边就没啥好处理的了
					if err1 != nil {
						log.Println(err1)
					}
				}
			}()
		}
	}()
}

// Uid 一般是从 jwt token 或者 session 里面取出来
// 这里模拟从 header 里面读取出来
func (s *WsGateway) Uid(req *http.Request) int64 {
	uidStr := req.Header.Get("uid")
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	return uid
}

func (s *WsGateway) Consume(msg *sarama.ConsumerMessage, event Event) error {
	// 我要消费
	conn, ok := s.conns.Load(event.Receiver)
	if !ok {
		log.Println("当前节点上没有这个用户，直接返回")
		return nil
	}
	return conn.Send(event.Msg)
}

type Conn struct {
	*websocket.Conn
}

func (c *Conn) Send(msg Message) error {
	val, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.Conn.WriteMessage(websocket.TextMessage, val)
}

// Message 前后端交互的数据格式
type Message struct {
	// 前端的序列号
	Seq string `json:"seq"`
	// 标记是什么类型的消息
	// 比如说图片，视频
	// {"type": "image", content:"http://myimage"}
	Type string `json:"type"`
	// 内容肯定有
	Content string `json:"content"`
	// 你发给谁？
	// channel id
	Cid int64
}
