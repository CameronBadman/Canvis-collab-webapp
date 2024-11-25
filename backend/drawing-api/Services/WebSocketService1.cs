using System;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace DrawingApi.Services
{
    public class WebSocketService1
    {
        public async Task HandleWebSocketAsync(WebSocket webSocket, CancellationToken cancellationToken)
        {
            var buffer = new byte[1024 * 4];
            Console.WriteLine("WebSocket connection established.");

            try
            {
                while (webSocket.State == WebSocketState.Open)
                {
                    var result = await webSocket.ReceiveAsync(new ArraySegment<byte>(buffer), cancellationToken);
                    Console.WriteLine($"Received message type: {result.MessageType}, EndOfMessage: {result.EndOfMessage}");

                    if (result.MessageType == WebSocketMessageType.Text)
                    {
                        var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
                        Console.WriteLine($"[Service1] Received message: {message}");

                        // Echo the message back to the client
                        var responseMessage = $"[Service1] Echo: {message}";
                        Console.WriteLine($"Sending message: {responseMessage}");
                        await webSocket.SendAsync(new ArraySegment<byte>(Encoding.UTF8.GetBytes(responseMessage)),
                            WebSocketMessageType.Text, result.EndOfMessage, cancellationToken);
                    }
                    else if (result.MessageType == WebSocketMessageType.Close)
                    {
                        Console.WriteLine("Received close request from client.");
                        if (webSocket.State == WebSocketState.Open)
                        {
                            await webSocket.CloseAsync(WebSocketCloseStatus.NormalClosure, "Closing connection", cancellationToken);
                            Console.WriteLine("WebSocket connection closed gracefully.");
                        }
                    }
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error handling WebSocket: {ex.Message}");
                if (webSocket.State != WebSocketState.Aborted)
                {
                    await webSocket.CloseAsync(WebSocketCloseStatus.InternalServerError, "Error occurred", cancellationToken);
                }
            }
        }
    }
}
