from helper import *
from data import *

async def join_room(websocket, data):
    for room in list_of_rooms:
        if room["room-id"] == int(data["room-id"]):
            if len(room["players"]) > 8:
                print("Room limit already reached")
                return 
            for player in room["players"]:
                if player["player-id"] == data["player-id"]:
                    await handle_update(websocket, room, player["name"], player["player-id"])
                    return
            add_player_to_room(websocket, data["player-id"], data["name"], room)
            await handle_update(websocket, room, data["name"], data["player-id"])
            return
    await error(websocket, "Couldn't Join Room")

async def create_room(websocket, data):
    number = generate_random_number()
    for room in list_of_rooms:
        if room["room-id"] == number:
            create_room(websocket, data)
            return
    list_of_rooms.append({
        "room-id": number,
        "room-size": 0,
        "room-buy-in": data["buy-in"],
        "room-big-blind-position": 2,
        "room-small-blind-position": 1,
        "room-dealer-position": 0,
        "room-pot": 0,
        "players": []
    })
    await send_room_number(websocket, number)

async def disconnect_from_room(data):
    player = find_player_in_list(data)
    if player:
        room = find_room_from_id(player["room-id"])
        for p in room["players"]:
            if player["player-id"] == p["player-id"]:
                room["players"].remove(p)
                room["room-size"] = room["room-size"] - 1
                await update_room_for_all_players(room)

async def check_for_player_metadata(websocket, player_data):
    player = find_player_in_list(player_data)
    if player:
        player_data["room-id"] = player["room-id"]
        await join_room(websocket, player_data)