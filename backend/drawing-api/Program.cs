var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.UseWebSockets();

app.Map("/ws", async context =>
{
    if (context.WebSockets.IsWebSocketRequest)
    {
        var socket = await context.WebSockets.AcceptWebSocketAsync();
        await HandleWebSocketAsync(socket);
    }
    else
    {
        context.Response.StatusCode = 400;
    }
});

app.Run();