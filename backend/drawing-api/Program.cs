var builder = WebApplication.CreateBuilder(args);

// Register services
builder.Services.AddSingleton<RedisService>();
builder.Services.AddSingleton<DrawingService>();
builder.Services.AddSingleton<DrawingMessageHandler>();
builder.Services.AddSingleton<WebSocketHandler>();

var app = builder.Build();

// Enable WebSocket support
app.UseWebSockets();

// WebSocket endpoint
app.UseEndpoints(endpoints =>
{
    endpoints.MapGet("/ws", async context =>
    {
        if (context.WebSockets.IsWebSocketRequest)
        {
            var webSocketHandler = context.RequestServices.GetRequiredService<WebSocketHandler>();
            var webSocket = await context.WebSockets.AcceptWebSocketAsync();
            await webSocketHandler.HandleAsync(context, webSocket);
        }
        else
        {
            context.Response.StatusCode = 400;
        }
    });
});

// Health check endpoint
app.MapControllers();

app.Run();
