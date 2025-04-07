// internal/ws/handler.go
package ws

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // cho demo, cần chỉnh kỹ khi production
}

type WSMessage struct {
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload"`
}

type CommandFunc func(msg WSMessage) map[string]interface{}

var commandRegistry = map[string]CommandFunc{}

func Register(name string, fn CommandFunc) {
	commandRegistry[name] = fn
}

func AutoRegister(fn CommandFunc) {
	name := getFunctionName(fn)
	Register(name, fn)
}

func Dispatch(msg WSMessage) map[string]interface{} {
	fn, ok := commandRegistry[msg.Command]
	if !ok {
		return map[string]interface{}{
			"status": "error",
			"error":  fmt.Sprintf("unknown command: %s", msg.Command),
		}
	}
	return fn(msg)
}

func getFunctionName(i interface{}) string {
	full := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(full, ".")
	return parts[len(parts)-1] // lấy tên hàm cuối cùng
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		response := Dispatch(msg)
		if err := conn.WriteJSON(response); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

func Connect() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", HandleWS)
	addr := ":" + os.Getenv("WS_PORT")
	fmt.Println("🚀 WebSocket server đang chạy tại", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("❌ Lỗi khi chạy server:", err)
	}
}
