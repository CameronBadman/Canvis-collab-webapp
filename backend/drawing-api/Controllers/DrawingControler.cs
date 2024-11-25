using Microsoft.AspNetCore.Mvc;

namespace DrawingApi.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class DrawingController : ControllerBase
    {
        [HttpGet("status")]
        public IActionResult GetStatus()
        {
            return Ok("Drawing API is running");
        }
    }
}
