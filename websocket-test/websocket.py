import asyncio
import websockets

async def test_websocket():
    # The WebSocket URL you are using, make sure it matches the endpoint in your .NET application
    uri = "ws://localhost:8001/service1"  # Adjust with your endpoint if needed

    # Connect to the WebSocket server
    async with websockets.connect(uri) as websocket:
        # Send a message to the server
        message = "Hello from Python WebSocket!"
        await websocket.send(message)
        print(f"Sent message: {message}")

        # Wait for a response from the server
        response = await websocket.recv()
        print(f"Received response: {response}")

# Run the WebSocket client
asyncio.run(test_websocket())
