package core

import (
	"../api"

	"github.com/Balhazraell/logger"
	"github.com/streadway/amqp"
)

// Server - singletone для работы с care частью сервера.GameServer
var Server server

type server struct {
	RoomIDByClient map[int]int
	Rooms          map[int]string

	//! Это список пользователей в очереди на подключение.
	// мы должны переодически проверять его на наличие там пользователей на подключение!
	PendingUsers []int

	shutdownLoop chan bool

	//--- RabbitMQ
	APIandCallbackMetods map[string]func(string, int)
	connectRMQ           *amqp.Connection
	channelRMQ           *amqp.Channel
}

// ServerStart - метод запуска игрового сервера.
func ServerStart() {
	Server = server{
		RoomIDByClient:       make(map[int]int),
		Rooms:                make(map[int]string),
		shutdownLoop:         make(chan bool),
		APIandCallbackMetods: fillMetods(),
	}

	go Server.loop()

	StartRabbitMQ()
}

// Stop - метод завершения работы игрового сервера.
func (serv *server) Stop() {
	serv.shutdownLoop <- true
}

func (serv *server) loop() {
	defer func() {
		serv.connectRMQ.Close()
		serv.channelRMQ.Close()
		logger.InfoPrint("Cервер закончил свою работу.")
	}()

	logger.InfoPrint("Cервер запущен.")

	for {
		select {
		case <-serv.shutdownLoop:
			return
		case clietnID := <-api.API.ClientConnectionChl:
			clientConnect(clietnID)
		case clientID := <-api.API.ClientDisconnectChl:
			clientDisconnect(clientID)
		case chunckStateData := <-api.API.SetChunckStateChl:
			setChunckState(
				chunckStateData.ClientID,
				chunckStateData.ChuncID,
			)
		case changeRoomStructureData := <-api.API.ChangeRoomChl:
			changeRoom(
				changeRoomStructureData.ClientID,
				changeRoomStructureData.RoomID,
			)
		}
	}
}

func fillMetods() map[string]func(string, int) {
	result := APIMetods
	for key, value := range CallbackMetods {
		result[key] = value
	}

	return result
}
