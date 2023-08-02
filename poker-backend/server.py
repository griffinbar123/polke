import asyncio
import websockets
import json
from handler import *

async def handler(websocket):
    while True:
        async for message in websocket:
            data = json.loads(message)  
            print(data["type"])
            if data["type"] == "player-metadata":
                await check_for_player_metadata(websocket, data["payload"])
            if data["type"] == "join-room":
                await join_room(websocket, data["payload"])
            if data["type"] == "create-room":
                await create_room(websocket, data["payload"])
            if data["type"] == "disconnect":
                await disconnect_from_room(data["payload"])


async def main():
    async with websockets.serve(handler, "", 8001):
        await asyncio.Future()  # run forever


if __name__ == "__main__":
    asyncio.run(main()) 


