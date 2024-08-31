using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Microsoft.Win32;
using OsuParsers.Beatmaps;
using OsuParsers.Database;
using OsuParsers.Database.Objects;
using shitweb;
using System.Collections;
using System.Net;
using System.Text.RegularExpressions;

public class Osudb
{
    private static Osudb instance = null;
    private static readonly object padlock = new object();
    public static string osufolder { get; private set; }
    public static CollectionDatabase CollectionDb { get; private set; }

    static Osudb()
    {
        var key = Registry.GetValue(
    @"HKEY_LOCAL_MACHINE\SOFTWARE\Classes\osu\shell\open\command",
    "",
    null
    );

        if (key != null)
        {
            string[] keyparts = key.ToString().Split('"');

            osufolder = Path.GetDirectoryName(keyparts[1]);

            Parse(osufolder);
        }
        else throw new Exception("Osu not Installed... ");
    }

    static void Parse(string filepath)
    {
        OsuDatabase osuDatabase = null;
        string file = "/osu!.db";
        if (File.Exists(filepath + file))
        {
            using (FileStream fileStream = new FileStream(filepath + file, FileMode.Open, FileAccess.Read, FileShare.ReadWrite, bufferSize: 4096, useAsync: true))
            {
                Console.WriteLine($"Parsing {file}");
                osuDatabase = OsuParsers.Decoders.DatabaseDecoder.DecodeOsu($"{filepath}{file}");
                Console.WriteLine($"Parsed {file}");

                fileStream.Close();
            }
        }

        file = "/collection.db";
        if (File.Exists(filepath + file))
        {
            using (FileStream fileStream = new FileStream(filepath + file, FileMode.Open, FileAccess.Read, FileShare.ReadWrite, bufferSize: 4096, useAsync: true))
            {
                Console.WriteLine($"Parsing {file}");
                CollectionDb = OsuParsers.Decoders.DatabaseDecoder.DecodeCollection($"{filepath}{file}");
                Console.WriteLine($"Parsed {file}");
            }
        }

        SqliteDB.Instance().setup(osuDatabase);
        osuDatabase = null;
        GC.Collect();
    }

    public static Osudb Instance
    {
        get
        {
            lock (padlock)
            {
                if (instance == null)
                {
                    instance = new Osudb();
                }
                return instance;
            }
        }
    }

    public static OsuParsers.Database.Objects.Collection GetCollectionbyName(string name) {

        return CollectionDb.Collections.FirstOrDefault(collection => collection.Name == name);
    }

    public static OsuParsers.Database.Objects.Collection GetCollectionbyIndex(int index)
    {

        return CollectionDb.Collections.ElementAtOrDefault(index);
    }

    public Collection GetCollection(int index) {

        var collection = GetCollectionbyIndex(index);
        if (collection == null) { return null; }

        List<Song> songs = new List<Song>();
        var activeId = "";

        collection.MD5Hashes.ForEach(hash =>
        {
            var beatmap = SqliteDB.GetSongByHash(hash);
            if (beatmap == null) { return; }

            songs.Add(beatmap);

        });

    return new Collection(collection.Name, songs.Count, songs);
    }

    public List<CollectionPreview> GetCollections(int limit, int offset)
    {

        List<CollectionPreview> collections = new List<CollectionPreview>();

        for (int i = offset; i < CollectionDb.Collections.Count - 1 && i < offset + limit; i++) {
                var collection = CollectionDb.Collections[i];

                var beatmap = SqliteDB.GetSongByHash(collection.MD5Hashes.FirstOrDefault());

                collections.Add(new CollectionPreview(index: i, name: collection.Name, previewimage: beatmap.previewimage, length: collection.Count));
            };

        return collections;
    }

    public static string getBG(string songfolder, string diff)
    {
        string folderpath = Path.Combine(songfolder, diff);
        string filepath = Path.Combine(osufolder, "Songs", folderpath);

        if (File.Exists(filepath))
        {
            string fileContents = File.ReadAllText($@"{filepath}"); // Read the contents of the file

            string pattern = @"\d+,\d+,""(?<image_filename>[^""]+\.[a-zA-Z]+)"",\d+,\d+";

            Match match = Regex.Match(fileContents, pattern);

            if (match.Success)
            {
                string background = Uri.EscapeDataString(match.Groups["image_filename"].Value);

                return Path.Combine(Uri.EscapeDataString(songfolder), background);
            }

            pattern = @"\d+,\d+,""(?<image_filename>[^""]+\.[a-zA-Z]+)""";

            match = Regex.Match(fileContents, pattern);

            if (match.Success)
            {
                string background = Uri.EscapeDataString(match.Groups["image_filename"].Value);
                return Path.Combine(Uri.EscapeDataString(songfolder), background);
            }
        }
        return "default-bg.png";
    }
}
