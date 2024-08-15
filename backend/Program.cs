using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.Hosting;
using System.Drawing;
using System.Drawing.Drawing2D;
using System.Drawing.Imaging;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddMemoryCache();

builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowAll",
        policy =>
        {
            policy.AllowAnyOrigin()
                  .AllowAnyHeader()
                  .AllowAnyMethod();
        });
});

var app = builder.Build();
app.UseCors("AllowAll");

app.MapGet("/ping", () => "pong");


// Define the API routes
app.MapGet("/api/v1/songs/{hash}", (string hash) =>
    {

        return Results.Ok(new { hash });
    });

app.MapGet("/api/v1/songs/recent", (int? limit, int? offset) =>
    {
        var limitValue = limit ?? 100; // default to 10 if not provided
        var offsetValue = offset ?? 0; // default to 0 if not provided

        return Results.Json(Osudb.Instance.GetRecent(limitValue, offsetValue));
    });

app.MapGet("/api/v1/songs/favorite", (int? limit, int? offset) =>
    {
        var limitValue = limit ?? 10; // default to 10 if not provided
        var offsetValue = offset ?? 0; // default to 0 if not provided

        return Results.Ok(new { Limit = limitValue, Offset = offsetValue, Message = "List of favorite songs" });
    });

app.MapGet("/api/v1/songs/{hash}", (string hash) =>
    {
        return Results.Ok($"Details for song with hash {hash}");
    });

app.MapGet("/api/v1/collections/", async (int? limit, int? offset, [FromServices] IMemoryCache cache) =>
    {
        const string cacheKey = "collections";

        if (!cache.TryGetValue(cacheKey, out var collections))
        {
         
            collections = Osudb.Instance.GetCollections();

            var cacheEntryOptions = new MemoryCacheEntryOptions()
                .SetSlidingExpiration(TimeSpan.FromDays(1)) 
                .SetAbsoluteExpiration(TimeSpan.FromDays(3));

            cache.Set(cacheKey, collections, cacheEntryOptions);
        }

        return Results.Json(collections);

    });

app.MapGet("/api/v1/collection/{index}", (int index) =>
    {
        return Results.Json(Osudb.Instance.GetCollection(index));
    });

app.MapGet("/api/v1/audio/{*fileName}", async (string fileName, HttpContext context) =>
{
    var decodedFileName = Uri.UnescapeDataString(fileName);
    var filePath = Path.Combine(Osudb.osufolder, "Songs", decodedFileName);

    if (!System.IO.File.Exists(filePath))
    {
        Console.WriteLine($"Not Found: {filePath}");
        return Results.NotFound();
    }

    var fileExtension = Path.GetExtension(filePath).ToLowerInvariant();
    var contentType = fileExtension switch
    {
        ".mp3" => "audio/mpeg",
        ".wav" => "audio/wav",
        ".ogg" => "audio/ogg",
        _ => "application/octet-stream",
    };

    var fileStream = new FileStream(filePath, FileMode.Open, FileAccess.Read, FileShare.Read);

    return Results.Stream(fileStream, contentType, enableRangeProcessing: true);
});


app.MapGet("/api/v1/images/{*filename}", async (string filename, int? h, int? w) =>
{
    var decodedFileName = Uri.UnescapeDataString(filename);
    var filePath = Path.Combine(Osudb.osufolder, "Songs", decodedFileName);

    if (!System.IO.File.Exists(filePath))
    {
        Console.WriteLine($"Not Found: {filePath}");
        return Results.NotFound();
    }

    var fileExtension = Path.GetExtension(filePath).ToLowerInvariant();
    var contentType = fileExtension switch
    {
        ".jpg" or ".jpeg" => "image/jpeg",
        ".png" => "image/png",
        ".gif" => "image/gif",
        ".bmp" => "image/bmp",
        ".webp" => "image/webp",
        _ => "application/octet-stream",
    };
    
    if (w == null || h == null)
    {
        var fileStream = new FileStream(filePath, FileMode.Open, FileAccess.Read, FileShare.Read, 4096, useAsync: true);
        return Results.Stream(fileStream, contentType, filename);
    }
    using var originalImage = new Bitmap(filePath);


    // If resizing is requested, resize the image
    Bitmap resizedImage;
    if (w.HasValue || h.HasValue)
    {
        resizedImage = ResizeImage(originalImage, w, h);
    }
    else
    {
        resizedImage = new Bitmap(originalImage); // Keep original size
    }

    // Convert the resized image to a memory stream
    var memoryStream = new MemoryStream();
    resizedImage.Save(memoryStream, GetImageFormat(fileExtension));
    memoryStream.Position = 0; // Reset stream position

    return Results.File(memoryStream, contentType);

});

static Bitmap ResizeImage(Image originalImage, int? width, int? height)
{
    int newWidth = width ?? originalImage.Width;
    int newHeight = height ?? originalImage.Height;

    if (width == null)
    {
        newWidth = originalImage.Width * newHeight / originalImage.Height;
    }
    else if (height == null)
    {
        newHeight = originalImage.Height * newWidth / originalImage.Width;
    }

    var resizedImage = new Bitmap(newWidth, newHeight);
    using (var graphics = Graphics.FromImage(resizedImage))
    {
        graphics.CompositingQuality = CompositingQuality.HighQuality;
        graphics.InterpolationMode = InterpolationMode.HighQualityBicubic;
        graphics.SmoothingMode = SmoothingMode.HighQuality;
        graphics.DrawImage(originalImage, 0, 0, newWidth, newHeight);
    }

    return resizedImage;
}

static ImageFormat GetImageFormat(string extension)
{
    return extension switch
    {
        ".jpg" or ".jpeg" => ImageFormat.Jpeg,
        ".png" => ImageFormat.Png,
        ".gif" => ImageFormat.Gif,
        ".bmp" => ImageFormat.Bmp,
        ".webp" => ImageFormat.Webp,
        _ => ImageFormat.Png,
    };
}

app.Run();
