package main

import (
	"log"
	"net/http"

	// 웹소켓 패키지
	"github.com/gorilla/websocket"
)

// make() 함수 : 초기화
// clients : 웹소켓 포인터로 초기화된 클라이언트 변수
var clients = make(map[*websocket.Conn]bool)

// broadcast : Message 타입의 브로드캐스트 채널 (통로) 생성. 이 채널을 통해 클라이언트에서 보낸 메시지를 보유 (큐잉) 하는 변수
var broadcast = make(chan Message)

// upgrader : HTTP 연결을 받고 웹소켓으로 업그레이드하는 기능 변수
var upgrader = websocket.Upgrader{}

// 메시지 구조체
type Message struct {
	Email    string `json:"email" `
	Username string `json:"username" `
	Message  string `json:"message" `
}

// 연결 핸들러
/**
i) GET 요청을 웹소켓으로 업그레이드
ii) 받은 요청을 클라이언트로 등록
iii) 웹소켓에서 메시지 수신 대기
iv) 메시지를 수신하면 이를 브로드캐스트 채널로 송신
**/
// nil 키워드 : Golang에서의 null
// := 연산자 : 변수 선언 시, var 키워드를 생략하고 변수명만으로 선언 가능. 단, func 내에서만 사용 가능
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// defer > finally 처럼 마지막에 Clean Up 작업
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error : %v", err)
			delete(clients, ws)
			break
		}
		// <- 연산자를 통해 브로드캐스트 채널로 전송
		broadcast <- msg
	}
}

// 브로드캐스트 채널에 큐잉된 메시지를 꺼냄
// 커넥션 핸들러에서 처리하고 있는 클라이언트 중 하나가 메시지를 보내면 이를 꺼내어 현재 연결된 모든 클라이언트에게 메시지를 보냄
func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error : %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

/**
go 키워드 : 고루틴 (GoRoutine). GO 런타임이 관리하는 경량 / 논리적 쓰레드.
이 키워드를 사용하여 함수를 호출하면, 런타임시 새로운 고루틴을 실행함.
고루틴은 비동기적으로 함수루틴을 실행하므로 여러 코드를 동시에 실행하는데 사용.
> 브로드캐스트 채널에 수신된 메시지를 꺼낼 때에는 비동기적 호출 필요
**/
func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	log.Println("HTTP Server started on Port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServce : ", err)
	}
}
