using System;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace DrawingApi.Services
{
    public class WebSocketService
    {
        public async Task HandleWebSocketAsync(WebSocket webSocket, CancellationToken cancellationToken)
        {
            var buffer = new byte[1024 * 4];
            Console.WriteLine("WebSocket connection established.");

            try
            {
                while (webSocket.State == WebSocketState.Open)
                {
                    // Receive WebSocket message
                    var result = await webSocket.ReceiveAsync(new ArraySegment<byte>(buffer), cancellationToken);
                    Console.WriteLine($"Received message type: {result.MessageType}, EndOfMessage: {result.EndOfMessage}");

                    if (result.MessageType == WebSocketMessageType.Text)
                    {
                        // Convert byte buffer to string
                        var message = Encoding.UTF8.GetString(buffer, 0, result.Count);

                        // Log the message and additional info
                        Console.WriteLine($"[Service] Received message: {message}");

                        // You can log more details such as:
                        // - Client address or request headers
                        // - Size of the message received
                        // - Time of the message
                        Console.WriteLine($"Received at: {DateTime.UtcNow}");
                        Console.WriteLine($"Message size: {result.Count} bytes");
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
                try
                {
                    if (webSocket.State != WebSocketState.Aborted)
                    {
                        await webSocket.CloseAsync(WebSocketCloseStatus.InternalServerError, "Error occurred", cancellationToken);
                        Console.WriteLine("WebSocket connection closed due to error.");
                    }
                }
                catch (Exception closeEx)
                {
                    Console.WriteLine($"Error closing WebSocket: {closeEx.Message}");
                }
            }
        }
    }
}
