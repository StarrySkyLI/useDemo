package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type UserMessage struct {
	UserID string
	Msg    string
}

type NewClient struct {
	UserID string
	Chan   chan string
}

type Event struct {
	Message       chan UserMessage
	NewClients    chan NewClient
	ClosedClients chan chan string
	TotalClients  map[string]map[chan string]bool
	clientToUser  map[chan string]string
}

type ClientChan chan string

func main() {
	router := gin.Default()

	stream := NewServer()

	// 模拟定时向特定用户发送消息
	go func() {
		for {
			time.Sleep(time.Second * 1)
			now := time.Now().Format("2006-01-02 15:04:05")
			currentTime := fmt.Sprintf("The Current Time Is %v", now)
			stream.Message <- UserMessage{
				UserID: "admin",
				Msg:    fmt.Sprintf("Admin message at %v", currentTime),
			}
			stream.Message <- UserMessage{
				UserID: "user1",
				Msg:    fmt.Sprintf("User message at %v", currentTime),
			}

		}
	}()

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "admin123",
		"user1": "password1",
	}))

	authorized.GET("/stream", HeadersMiddleware(), stream.serveHTTP(), func(c *gin.Context) {
		v, ok := c.Get("clientChan")
		if !ok {
			return
		}
		clientChan, ok := v.(ClientChan)
		if !ok {
			return
		}
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-clientChan; ok {
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
		Message:       make(chan UserMessage),
		NewClients:    make(chan NewClient),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[string]map[chan string]bool),
		clientToUser:  make(map[chan string]string),
	}
	go event.listen()
	return event
}

func (stream *Event) listen() {
	for {
		select {
		case client := <-stream.NewClients:
			userID := client.UserID
			clientChan := client.Chan

			if _, exists := stream.TotalClients[userID]; !exists {
				stream.TotalClients[userID] = make(map[chan string]bool)
			}
			stream.TotalClients[userID][clientChan] = true
			stream.clientToUser[clientChan] = userID
			log.Printf("%s connected. Total users: %d", userID, len(stream.TotalClients))

		case clientChan := <-stream.ClosedClients:
			if userID, exists := stream.clientToUser[clientChan]; exists {
				delete(stream.TotalClients[userID], clientChan)
				if len(stream.TotalClients[userID]) == 0 {
					delete(stream.TotalClients, userID)
				}
				delete(stream.clientToUser, clientChan)
				close(clientChan)
				log.Printf("%s disconnected. Remaining users: %d", userID, len(stream.TotalClients))
			}

		case msg := <-stream.Message:
			if clients, ok := stream.TotalClients[msg.UserID]; ok {
				for clientChan := range clients {
					select {
					case clientChan <- msg.Msg:
					default:
						log.Printf("Client %s channel full, message dropped", msg.UserID)
					}
				}
			}
		}
	}
}

func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet(gin.AuthUserKey).(string)
		clientChan := make(ClientChan)

		stream.NewClients <- NewClient{
			UserID: userID,
			Chan:   clientChan,
		}

		go func() {
			<-c.Done()
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)
		c.Next()
	}
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Next()
	}
}
