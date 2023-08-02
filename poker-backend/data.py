list_of_rooms = []
list_of_players = []

def find_room_from_player(player):
    for room in list_of_rooms:
        for p in room["players"]:
            if player["player-id"] == p["player-id"]:
                return room["room-id"]
    print("the player is not in a room")
    return None

def find_player_in_list(player):
    for p in list_of_players:
        if p["player-id"] == player["player-id"]:
            return p
    print("no player with that id")
        
def find_room_from_id(room_id):
    for room in list_of_rooms:
        if room["room-id"] == room_id:
            return room
    print("no room with that id")

def remove_from_players_set(player):
    for p in list_of_players:
        if p["player-id"] == player["player-id"]: 
            list_of_players.remove(p)

def add_to_players_set(player):
    list_of_players.append(player)