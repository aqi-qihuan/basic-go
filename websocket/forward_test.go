package websocket

import (
	"log"
	"net/http"
	"testing"

	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
)

type Hub struct {
	conns *syncx.Map[string, *websocket.Conn]
}

func (h *Hub) AddConn(name string, conn *websocket.Conn) {
	h.conns.Store(name, conn)
	go func() {
		for {
			// 接收数据
			typ, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("接收 websocket 数据失败", err)
				return
			}
			// 转发数据
			// 你的返回值决定了要不要继续遍历
			h.conns.Range(func(key string, value *websocket.Conn) bool {
				if key == name {
					// 我自己就不需要转发了
					return true
				}
				err1 := value.WriteMessage(typ, msg)
				if err1 != nil {
					log.Println(err)
					// 记录日志
				}
				return true
			})
		}
	}()
}

func TestHub(t *testing.T) {
	upgrader := websocket.Upgrader{}
	hub := &Hub{conns: &syncx.Map[string, *websocket.Conn]{}}
	// 我们假定，websocket 请求发到这里
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// responseHeader 可以不传
		name := request.URL.Query().Get("name")
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			writer.Write([]byte("初始化 websocket 失败"))
			return
		}
		hub.AddConn(name, conn)
	})
	http.ListenAndServe(":8081", nil)
}
