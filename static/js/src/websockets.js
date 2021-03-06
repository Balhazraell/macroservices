'use strict';

var main = require('./main');

var ws

// Набор функций получаемых от сервера
var handlers = {
    'set_grid': set_grid,
    'send_error': send_error,
    'set_rooms_catalog': set_rooms_catalog,
    'set_select_room': set_select_room,
};

// Пошла работа с websockets
function connect(){
    // Здесь происходит падение если не можем подключится, надо красиво обработать...
    ws = new WebSocket('ws://127.0.0.1:8081/appgame');
    ws.onopen = open;
    ws.onclose = close;
    ws.onmessage = message;
}

// websocket стартанул.
function open(event){
    console.log('websocket is open!');
}

// websocket закрылся.
function close(event){
    console.log('websocket is close!');
}

// пришло сообщение по websocket.
function message(event){
    var data = JSON.parse(event.data);
    handlers[data['handler_name']](JSON.parse(data['data']));
}

// ------------- incoming ------------------
// Пришла сетка.
function set_grid(new_map){
    main.set_grid(new_map);
}

// Пришла ошибка от сервера.
function send_error(message){
    main.send_error(message);
}

// Отправляем запрос на постановку символа в чанк
function set_chunck_state(chunck_id){
    var data = {
        'chunck_id': chunck_id,
    }

    var message = {
        'handler_name': 'setChunckState',
        'data': JSON.stringify(data),
    }

    ws.send(JSON.stringify(message));
}

function sendChangeRoomID(roomID){
    var data = {
        'room_id': parseInt(roomID),
    }

    var message = {
        'handler_name': 'chengeRoomID',
        'data': JSON.stringify(data),
    }

    console.log("Отправляем id = ", message);

    ws.send(JSON.stringify(message));
}

function set_rooms_catalog(roomsIDs){
    console.log("Католог комнат:", roomsIDs);
    main.setRoomCatalog(roomsIDs);
}

function set_select_room(roomID){
    console.log("Задана комната с id:", roomID);
    main.setSelectRoom(roomID);
}

exports.connect = connect;
exports.set_chunck_state = set_chunck_state;
exports.sendChangeRoomID = sendChangeRoomID;