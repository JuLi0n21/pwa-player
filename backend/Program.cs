using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.Hosting;
using shitweb;
using System;
using System.Diagnostics;
using System.Drawing;
using System.Drawing.Drawing2D;
using System.Drawing.Imaging;
using System.Drawing.Text;
using System.Net.Http;
using System.Text.RegularExpressions;

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

var apiClient = new ApiClient();


app.MapGet("/ping", () => "pong");


app.MapGet("/login", () => {
    return Results.Redirect("https://proxy.illegalesachen.download/login");
});

app.MapGet("/api/v1/songs/{hash}", (string hash) =>
{

    return Results.Ok(new { hash });
});

app.MapGet("/api/v1/songs/recent", (HttpContext httpContext, int? limit, int? offset) =>
{
        var limitValue = limit ?? 100; // default to 10 if not provided
        var offsetValue = offset ?? 0; // default to 0 if not provided

        return Results.Json(SqliteDB.GetByRecent(limitValue, offsetValue));
    });

app.MapGet("/api/v1/songs/favorite", (int? limit, int? offset) =>
    {
        var limitValue = limit ?? 100; // default to 10 if not provided
        var offsetValue = offset ?? 0; // default to 0 if not provided

        return Results.Ok(new { Limit = limitValue, Offset = offsetValue, Message = "List of favorite songs" });
    });

app.MapGet("/api/v1/songs/{hash}", (string hash) =>
    {
        return Results.Ok($"Details for song with hash {hash}");
    });

app.MapGet("/api/v1/collections/", async (int? limit, int? offset, [FromServices] IMemoryCache cache) =>
    {

        var limitValue = limit ?? 100; // default to 10 if not provided
        var offsetValue = offset ?? 0;

        string cacheKey = $"collections_{offsetValue}_{limitValue}";

        if (!cache.TryGetValue(cacheKey, out var collections))
        {
         
            collections = Osudb.Instance.GetCollections(limit: limitValue, offset: offsetValue);

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

    if (!File.Exists(filePath))
    {
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

app.MapGet("/api/v1/search/active", async (string? q) =>
{
    return Results.Ok(SqliteDB.activeSearch(q));
});

app.MapGet("/api/v1/search/artist", async (string? q, int? limit, int? offset) =>
{
    var limitv = limit ?? 100;
    var offsetv = offset ?? 0;

    return Results.Ok(SqliteDB.GetArtistSearch(q, limitv, offsetv));
});

app.MapGet("/api/v1/search/songs", async (string? q, int? limit, int? offset) =>
{
    var limitv = limit ?? 100;
    var offsetv = offset ?? 0;
    return Results.Ok();
});

app.MapGet("/api/v1/search/collections", async (string? q, int? limit, int? offset) =>
{
    var limitv = limit ?? 100;
    var offsetv = offset ?? 0;
    return Results.Ok();
});

app.MapGet("/api/v1/images/{*filename}", async (string filename, int? h, int? w) =>
{
    var decodedFileName = Uri.UnescapeDataString(filename);
    var filePath = Path.Combine(Osudb.osufolder, "Songs", decodedFileName);

    if (!File.Exists(filePath))
    {
        filePath = "default-bg.png";
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

Osudb.Instance.ToString();
startCloudflared();

Task.Run(() =>
{
    Thread.Sleep(500);
    if (!apiClient.LoadCookies())
    {
        Console.WriteLine("Please visit this link and paste the Value back into here: ");

        var cookie = Console.ReadLine();

        apiClient.SaveCookies(cookie);
    }

    Console.WriteLine("Ur Osu songs should now be available, please delete the cookies.json if it doesnt show up and try again.");
});

await apiClient.InitializeAsync();
app.Run();

async Task startCloudflared() {
   
    var process = new Process
    {
        StartInfo = new ProcessStartInfo
        {
            FileName = "cloudflared",
            Arguments = "tunnel --url http://localhost:5153",
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
            CreateNoWindow = true
        }
    };

    process.ErrorDataReceived += (sender, e) => {
        if (!string.IsNullOrEmpty(e.Data))
        {
            ParseForUrls(e.Data);
        }
    };

    process.Start();

    process.BeginErrorReadLine();

    await Task.Run(() => process.WaitForExit());
}

void ParseForUrls(string data)
{
    var urlRegex = new Regex(@"https?://[^\s]*\.trycloudflare\.com");
    var matches = urlRegex.Matches(data);

    foreach (Match match in matches)
    {
        Console.WriteLine($"Login here if not done already: {match.Value}/login");
        apiClient.UpdateSettingsAsync(match.Value);
    }
}
