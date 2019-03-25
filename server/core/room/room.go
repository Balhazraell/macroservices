package room

import (
	"../../logger"
)

var Room RoomStruct

// Chunc описывает струтуру участка игрового пространства.
type Chunc struct {
	ID          int      `json:"id"`
	State       int      `json:"state"`
	Сoordinates [][2]int `json:"coordinates"`
}

type RoomStruct struct {
	ID      int
	Map     map[int]*Chunc
	clients []int

	// Переменные логики.
	GameState int // Делаем крестики нолики, по этому 2 состояния - ходит один потом другой.

	// Каналы
	shutdownLoop chan bool
	updateMap    chan bool
}

// StartNewRoom - метод запуска новой комнаты.
// На вход подается id комнаты котурую надо создать.
func StartNewRoom(id int) {
	Room := RoomStruct{
		Map:          make(map[int]*Chunc),
		shutdownLoop: make(chan bool),
		updateMap:    make(chan bool),
	}

	createMap()

	go Room.loop()
}

// Stop - Останавлием работу комнаты
func (room *RoomStruct) Stop() {
	// ...какая-нибудь логика завершения работы.
	room.shutdownLoop <- true
}

func (room *RoomStruct) loop() {
	defer func() {
		logger.InfoPrintf("Комната с id=%v закончила работу.", room.ID)
	}()

	logger.InfoPrintf("Комната с id=%v начала работу.", room.ID)

	for {
		// Обновление логики происходит тут.

		select {
		case <-room.shutdownLoop:
			return

		// Даже не знаю на сколько целесообразно делать это в отдельном потоке.
		// Мсль была в том, что update карт должен произоти не моментально после изменений
		// но хз на сколько это грамотоное решение.
		case <-room.updateMap:
			updateClientsMap(Room.clients)
		}
	}
}
