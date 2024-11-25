using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using DrawingApi.Services;
using System.Net.WebSockets;


var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddSingleton<WebSocketService1>();
builder.Services.AddSingleton<WebSocketService2>();


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
        var webSocket = await context.WebSockets.AcceptWebSocketAsync();

        Console.WriteLine($"WebSocket request path: {context.Request.Path}");

        if (context.Request.Path == "/service1")
        {
            Console.WriteLine("Routing to WebSocketService1");
            var webSocketService1 = app.Services.GetRequiredService<WebSocketService1>();
            await webSocketService1.HandleWebSocketAsync(webSocket, context.RequestAborted);
        }
        else if (context.Request.Path == "/service2")
        {
            Console.WriteLine("Routing to WebSocketService2");
            var webSocketService2 = app.Services.GetRequiredService<WebSocketService2>();
            await webSocketService2.HandleWebSocketAsync(webSocket, context.RequestAborted);
        }
        else
        {
            Console.WriteLine("Unknown WebSocket request path.");
            await webSocket.CloseAsync(WebSocketCloseStatus.EndpointUnavailable, "Invalid WebSocket path", context.RequestAborted);
        }
    }
    else
    {
        await next();
    }
});



app.Run();
