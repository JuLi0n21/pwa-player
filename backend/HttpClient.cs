using OsuParsers.Enums.Replays;
using System;
using System.Diagnostics.CodeAnalysis;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Text.Json;
using System.Threading.Tasks;

namespace shitweb
{
    public class ApiClient
    {
        private const string CookiesFilePath = "cookies.json";
        private HttpClient _client;
        private HttpClientHandler _handler;

        public ApiClient()
        {
            _handler = new HttpClientHandler();
            _client = new HttpClient(_handler);
        }

        public async Task InitializeAsync()
        {
            LoadCookies();

        }

        public void SaveCookies(String cookie)
        {
            var cookies = _handler.CookieContainer.GetCookies(new Uri("https://proxy.illegalesachen.download"));
            var Cookie = new CookieInfo();

            Cookie.Name = "session_cookie";
            Cookie.Value = cookie;
            var json = JsonSerializer.Serialize(Cookie, new JsonSerializerOptions { WriteIndented = true });
            File.WriteAllText(CookiesFilePath, json);
            LoadCookies();
        }

        public Boolean LoadCookies()
        {
            if (File.Exists(CookiesFilePath))
            {
                var json = File.ReadAllText(CookiesFilePath);
                var cookieInfo = JsonSerializer.Deserialize<CookieInfo>(json);

                var cookie = new Cookie(cookieInfo.Name, cookieInfo.Value);
                _handler.CookieContainer.Add(new Uri("https://proxy.illegalesachen.download"), cookie);
                return true;
            }

            return false;
        }
    

        public async Task UpdateSettingsAsync(string endpoint)
        {
            var requestContent = new JsonContent(new { endpoint = endpoint });
            var response = await _client.PostAsync("https://proxy.illegalesachen.download/settings", requestContent);
            try
            {
                response.EnsureSuccessStatusCode();
            }
            catch (Exception ex) { 
                System.Console.WriteLine(ex.Message);
            }
        }
    }

    public class CookieInfo
    {
        public string Name { get; set; }
        public string Value { get; set; }
    }

    public class JsonContent : StringContent
    {
        public JsonContent(object obj)
            : base(JsonSerializer.Serialize(obj), System.Text.Encoding.UTF8, "application/json")
        {
        }
    }
}