using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using DrawingApi.Services;
using System.Net.WebSockets;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddSingleton<WebSocketService>();  // Ensure WebSocketService is registered.

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
        Console.WriteLine("Routing to WebSocketService");

        // Get the WebSocket service and handle the WebSocket request
        var webSocketService = app.Services.GetRequiredService<WebSocketService>();
        await webSocketService.HandleWebSocketAsync(webSocket, context.RequestAborted);
    }
    else
    {
        // Handle the case where the path isn't a WebSocket request
        Console.WriteLine("Unknown WebSocket request path.");

        // Send a bad request response (set status code to 400)
        context.Response.StatusCode = 400;
        await context.Response.WriteAsync("Invalid WebSocket path.");
    }

    // Call the next middleware in the pipeline
    await next.Invoke();
});


// Map HTTP controllers (if any)
app.MapControllers();

app.Run();
