using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using DrawingApi.Services;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddSingleton<WebSocketService>(); // Register WebSocket service

var app = builder.Build();

// Configure the HTTP request pipeline.
app.UseRouting();

app.MapControllers();

// Enable WebSocket support
app.UseWebSockets();
app.Use(async (context, next) =>
{
    if (context.WebSockets.IsWebSocketRequest)
    {
        using var webSocket = await context.WebSockets.AcceptWebSocketAsync();
        var websocketService = app.Services.GetRequiredService<WebSocketService>();
        await websocketService.HandleWebSocketAsync(webSocket, context.RequestAborted);
    }
    else
    {
        await next();
    }
});

app.Run();
