public class WebSocketHandler
{
    private readonly DrawingMessageHandler _messageHandler;

    public WebSocketHandler(DrawingMessageHandler messageHandler)
    {
        _messageHandler = messageHandler;
    }

    public async Task HandleAsync(HttpContext context, WebSocket webSocket)
    {
        var buffer = new byte[1024 * 4];
        var result = await webSocket.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);

        while (!result.CloseStatus.HasValue)
        {
            var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
            await _messageHandler.ProcessMessageAsync(message, webSocket);
            result = await webSocket.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        }

        await webSocket.CloseAsync(result.CloseStatus.Value, result.CloseStatusDescription, CancellationToken.None);
    }
}
