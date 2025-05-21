package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"useDemo/base-common/map_util"
)

type Connection struct {
	UserID int64
	Msg    chan []byte
}
type Message struct {
	UserID int64
	Msg    []byte
}

// 连接管理模块
type ConnManager struct {
	connections *map_util.ShardMap
}
type Event struct {
	MessageChan chan Message
	Manager     *ConnManager
	ClosedChan  chan *Connection
}

func NewManager() *ConnManager {

	return &ConnManager{
		connections: map_util.NewShardMap(),
	}
}

func main() {
	router := gin.Default()

	stream := NewServer()

	// 模拟定时向特定用户发送消息
	go func() {
		for {
			time.Sleep(time.Second * 1)
			stream.MessageChan <- Message{
				UserID: 1,
				Msg:    []byte(fmt.Sprintf("Admin message at %v", time.Now().Format("15:04:05"))),
			}
			stream.MessageChan <- Message{
				UserID: 2,
				Msg:    []byte(fmt.Sprintf("User message at %v", time.Now().Format("15:04:05"))),
			}
		}
	}()

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"1": "admin123",
		"2": "password1",
	}))

	authorized.GET("/stream", stream.serveHTTP(), func(c *gin.Context) {
		v, ok := c.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(*Connection)
		if !ok {
			return
		}
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-clientChan.Msg; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
	})

	router.StaticFile("/", "./go_demo/sse_demo/gin_demo_v0.1/public/index.html")
	router.Run(":8085")
}

func NewServer() *Event {
	event := &Event{
		Manager:     NewManager(),
		MessageChan: make(chan Message),
		ClosedChan:  make(chan *Connection),
	}

	go event.listen()
	return event
}

func (stream *Event) listen() {
	for {
		select {

		case message := <-stream.MessageChan:

			if v, exists := stream.Manager.connections.Load(strconv.FormatInt(message.UserID, 10)); exists {
				userConn := v.(*Connection)
				select {
				case userConn.Msg <- message.Msg:
				default:
					log.Printf("Client %d channel full, message dropped", userConn.UserID)
				}
			}

		case clientChan := <-stream.ClosedChan:

			if stream.Manager.connections.Delete(strconv.FormatInt(clientChan.UserID, 10)) {
				close(clientChan.Msg)
				log.Printf("Removed client. %d registered clients", clientChan.UserID)
			}

		}
	}
}

func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		ID := c.MustGet(gin.AuthUserKey).(string)
		userID, _ := strconv.ParseInt(ID, 10, 64)
		conn := &Connection{
			UserID: userID,
			Msg:    make(chan []byte),
		}
		stream.Manager.connections.Store(strconv.FormatInt(userID, 10), conn)
		log.Printf("Client %d connected", userID)

		go func() {
			<-c.Request.Context().Done()
			log.Printf("Client %d disconnected", userID)
			stream.ClosedChan <- conn
		}()

		c.Set("clientChan", conn)
		c.Next()
	}
}
