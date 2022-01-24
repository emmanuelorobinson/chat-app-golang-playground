package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		log.Print("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	router := gin.Default()

	//connect to socket.io
	//let router call server
	router.GET("/socket.io/", gin.WrapH(server))

	router.LoadHTMLGlob("public/*")

	//set asset to static
	router.Static("/assets", "assets")

	router.GET("/jsons", func(c *gin.Context) {
		//allow cross origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//http.Handle("/", http.FileServer(http.Dir("./public")))
	//log.Println("Serving at localhost:8000...")
	//log.Fatal(http.ListenAndServe(":8000", nil))
	router.Run(":8000")
}
