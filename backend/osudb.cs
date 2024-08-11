using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Microsoft.Win32;
using OsuParsers.Beatmaps;
using OsuParsers.Database;
using OsuParsers.Database.Objects;
using System.Collections;
using System.Net;
using System.Text.RegularExpressions;

public class Osudb
{
    private static Osudb instance = null;
    private static readonly object padlock = new object();
    public static string osufolder { get; private set; }
    public static OsuDatabase osuDatabase { get; private set; }
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
        string file = "/osu!.db";
        if (File.Exists(filepath + file))
        {
            using (FileStream fileStream = new FileStream(filepath + file, FileMode.Open, FileAccess.Read, FileShare.ReadWrite, bufferSize: 4096, useAsync: true))
            {
                Console.WriteLine($"Parsing {file}");
                osuDatabase = OsuParsers.Decoders.DatabaseDecoder.DecodeOsu($"{filepath}{file}");
                Console.WriteLine($"Parsed {file}");
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

    public DbBeatmap GetBeatmapbyHash(string Hash)
    {
        if (Hash == null || osuDatabase == null)
        {
            return null;
        }
        return osuDatabase.Beatmaps.FirstOrDefault(beatmap => beatmap.MD5Hash == Hash);
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
        var activeId = 0;

        collection.MD5Hashes.ForEach(hash =>
        {
            var beatmap = GetBeatmapbyHash(hash);
            if (beatmap == null) { return; }

            if (activeId == beatmap.BeatmapSetId) { return; }
            activeId = beatmap.BeatmapSetId;
            //todo
            string img = getBG(beatmap.FolderName, beatmap.FileName);

            songs.Add(new Song(hash: beatmap.MD5Hash, name: beatmap.Title, artist: beatmap.Artist, length: beatmap.TotalTime, url: $"{beatmap.FolderName}/{beatmap.AudioFileName}" , previewimage: img, mapper: beatmap.Creator));

        });

    return new Collection(collection.Name, songs.Count, songs);
    }

    public List<CollectionPreview> GetCollections()
    {

        List<CollectionPreview> collections = new List<CollectionPreview>();

        for (int i = 0; i < CollectionDb.Collections.Count; i++) {
                var collection = CollectionDb.Collections[i];

                var beatmap = GetBeatmapbyHash(collection.MD5Hashes.FirstOrDefault());

                //todo
                string img = getBG(beatmap.FolderName, beatmap.FileName);

                collections.Add(new CollectionPreview(index: i, name: collection.Name, previewimage: img, length: collection.Count));
            };

        return collections;
    }

    public List<Song> GetRecent(int limit, int offset)
    {
        var recent = new List<Song>();
        if(limit > 100 && limit < 0) {
            limit = 100;
        }

        var size = osuDatabase.Beatmaps.Count -1;
        for (int i = size - offset; i > size - offset - limit; i--)
        {
            var beatmap = osuDatabase.Beatmaps.ElementAt(i);
            if (beatmap == null) {
                continue;
            }

            string img = getBG(beatmap.FolderName, beatmap.FileName);

            recent.Add(new Song(
                name: beatmap.FileName, 
                hash: beatmap.MD5Hash, 
                artist: beatmap.Artist, 
                length: beatmap.TotalTime, 
                url: $"{beatmap.FolderName}/{beatmap.AudioFileName}",
                previewimage: img, 
                mapper: beatmap.Creator));
        }

        return recent;
    }

    public List<Song> GetFavorites()
    {
        var recent = new List<Song>();
        /*
        osuDatabase.Beatmaps.ForEach(beatmap =>
        {
            Console.WriteLine(beatmap.LastModifiedTime);
        });
        */
        return null;
    }

    private static string getBG(string songfolder, string diff)
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
                string background = match.Groups["image_filename"].Value;

                return Path.Combine(songfolder, background);
            }

            pattern = @"\d+,\d+,""(?<image_filename>[^""]+\.[a-zA-Z]+)""";

            match = Regex.Match(fileContents, pattern);

            if (match.Success)
            {
                string background = match.Groups["image_filename"].Value;
                return Path.Combine(songfolder, background);
            }
        }
        return null;
    }
}
