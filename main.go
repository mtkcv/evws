package main

import (
	"flag"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)

type taskWrapper struct {
	do func()
}

func runClient() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:12345", nil)
	if err != nil {
		logger.Fatalf("websocket connect error:%s", err.Error())
	}
	conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second))
	conn.WriteMessage(websocket.TextMessage, []byte("hello"))
	conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(4000, "close"),
		time.Now().Add(time.Second))
	conn.Close()
}

func main() {
	var r bool
	flag.BoolVar(&r, "r", false, "run client")
	flag.Parse()

	if r {
		LoggerInit(true)
		runClient()
		return
	}

	LoggerInit(false)
	poll := NewPoller()
	defer poll.Close()
	tp := NewTaskPool(256)
	defer tp.Close()

	conns := make(map[int]*websocket.Conn)
	go poll.Wait(func(fd, ev int) {
		conn := conns[fd]
		if conn == nil {
			logger.Errorf("fd %d not found websocket conn", fd)
			poll.Delete(fd)
			conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(4000, "read error"),
				time.Now().Add(time.Second))
			conn.Close()
			return
		}
		if ev&PollIn != 0 {
			messageType, data, err := conn.ReadFrame()
			if err != nil {
				logger.Errorf("fd %d(%s) read frame error, %s", fd, conn.RemoteAddr().String(), err.Error())
				poll.Delete(fd)
				conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(4000, "read error"),
					time.Now().Add(time.Second))
				conn.Close()
				return
			}
			if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
				logger.Infof("fd %d(%s) sent %d frame", fd, conn.RemoteAddr().String(), messageType)
				return
			}
			logger.Infof("fd %d(%s) sent %d %s", fd, conn.RemoteAddr().String(), messageType, string(data))
		}

		if ev&PollHup != 0 || ev&PollErr != 0 {
			poll.Delete(fd)
			conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(4000, "read error"),
				time.Now().Add(time.Second))
			conn.Close()
		}
	})

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		upg := websocket.Upgrader{
			HandshakeTimeout: time.Second * 5,
		}
		conn, err := upg.Upgrade(resp, req, nil)
		if err != nil {
			fmt.Println(err)
		}

		fd := websocketFD(conn)
		err = poll.Add(websocketFD(conn))
		if err != nil {
			fmt.Println(err)
		}
		conns[fd] = conn
	})

	err := http.ListenAndServe("127.0.0.1:12345", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func websocketFD(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
