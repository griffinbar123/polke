from data import *
import random
import json

def generate_random_number():
    return random.randint(1111,9999)



def format_room_state_to_send(room):
    r = {}
    r["room-id"] = room["room-id"]
    r["room-size"] = room["room-size"]
    r["room-buy-in"] = room["room-buy-in"]
    r["room-pot"] = room["room-pot"]
    r["room-big-blind-position"] = room["room-big-blind-position"]
    r["room-small-blind-position"] = room["room-small-blind-position"]
    r["room-dealer-position"] = room["room-dealer-position"]
    r["players"] = []
    for i in room["players"]:
        r["players"].append({
            "player-id" : i["player-id"],
            "name" : i["name"],
            "position" : i["position"],
            "cards" : i["cards"],
            "balance": i["balance"]
        })
    return r

def make_player(player_id, room_id=None, name='', connection=None, position=None, cards=None, balance=None, simple=False):
    if simple:
        return {
            "player-id": player_id,
            "room-id": room_id
        }
    return {
        "name": name,
        "player-id": player_id,
        "connection": connection,
        "position": position,
        "cards": cards,
        "balance": balance
    } 

def add_player_to_room(websocket, player_id, name, room):
    smallest_unfilled_position = 0
    while True:
        z = False
        for p in room["players"]:
            if p["position"] == smallest_unfilled_position:
                smallest_unfilled_position += 1
                z = True
                break
        if z is False:
            break
    room["players"].append(make_player(player_id, name=name, position=smallest_unfilled_position, cards=[0, 0], balance=room["room-buy-in"], connection=websocket))
    room["room-size"] = room["room-size"] + 1

async def update_room_for_all_players(room):
    for p in room["players"]:
        await send_room_state(p["connection"], room)

async def send_room_state(websocket, room):
    event = {
        "type": "room-state",
        "payload": format_room_state_to_send(room)
    }
    await send_event(websocket, event)

async def store_connection(websocket, name, player_id, room_id):
    event = {
        "type": "player-metadata", 
        "payload": {"name" : name, "player-id": player_id}
    }
    player = make_player(player_id, room_id=room_id, simple=True)
    remove_from_players_set(player)
    add_to_players_set(player)
    await send_event(websocket, event)

async def handle_update(websocket, room, name, player_id):
    await store_connection(websocket, name, player_id, room["room-id"])
    await update_room_for_all_players(room)

async def error(websocket, message):
    event = {
        "type": "error",
        "message": message,
    }
    await send_event(websocket, event)

async def send_event(websocket, event):
     await websocket.send(json.dumps(event))

async def send_room_number(websocket, room_number):
    event = {
        "type": "send-room-number", 
        "payload": {"room-number": room_number}
    }
    await send_event(websocket, event)
