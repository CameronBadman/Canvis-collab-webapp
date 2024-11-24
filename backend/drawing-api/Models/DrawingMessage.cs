
public class DrawingMessage
{
    public string CanvasId { get; set; }
    public string UserId { get; set; }
    public string Action { get; set; } // e.g., "draw", "clear", "move"
    public object Data { get; set; }  // Payload (e.g., drawing coordinates)
}
