public class DrawingService
{
    private readonly RedisService _redisService;

    public DrawingService(RedisService redisService)
    {
        _redisService = redisService;
    }

    public async Task UpdateCanvasAsync(string canvasId, DrawingMessage message)
    {
        // Update Redis with the new state
        await _redisService.SaveCanvasStateAsync(canvasId, message);
    }

    public async Task<string> GetCanvasStateAsync(string canvasId)
    {
        return await _redisService.GetCanvasStateAsync(canvasId);
    }
}
