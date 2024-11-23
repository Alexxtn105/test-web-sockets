package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// upgrader - главная задача - преобразование (апгрейд) http-соединения в WebSocket-соединение
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // размер буфера чтения
	WriteBufferSize: 1024, // размер буфера записи
	CheckOrigin: func(r *http.Request) bool { // разрешаем доступ всем
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Апгрейдим http-соединение до WebSocket-соединения
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//обязательно закрываем соединение по завершении работы
	defer conn.Close()

	// обработка входящих сообщений (используя метод readMessage соединения)
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		//выводим принятое сообщение
		fmt.Printf("Received message: %s\n", data)

		// отправляем клиенту то же самое сообщение обратно
		if err := conn.WriteMessage(messageType, data); err != nil {
			fmt.Println(err)
			return
		}
	}

}

func main() {
	//fmt.Println("Hello, World!")
	http.HandleFunc("/ws", handleWebSocket)
	addr := "localhost:8080"
	fmt.Printf("Server starting at: %s/ws\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}

}
