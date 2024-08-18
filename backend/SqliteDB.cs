using OsuParsers.Beatmaps;
using System.Collections.Generic;
using System.Data;
using System.Data.SqlClient;
using System.Data.SQLite;
using System.Reflection.PortableExecutable;

namespace shitweb
{
    public class SqliteDB
    {
        private static SqliteDB instance;
        private const string filename = "Beatmaps.db";
        private const string dburl = $"Data Source={filename};Version=3;";

        static SqliteDB() { }

        public static SqliteDB Instance()
        {
            if (instance == null)
            {

                instance = new SqliteDB();
            }
            return instance;
        }

        public void setup(OsuParsers.Database.OsuDatabase osuDatabase)
        {

            int count = 0;

            using (var connection = new SQLiteConnection(dburl))
            {
                connection.Open();

                string createBeatmaps = @"
                    CREATE TABLE IF NOT EXISTS Beatmap (
                        BeatmapId INTEGER DEFAULT 0,
                        Artist TEXT DEFAULT '?????',
                        ArtistUnicode TEXT DEFAULT '?????',
                        Title TEXT DEFAULT '???????',
                        TitleUnicode TEXT DEFAULT '???????',
                        Creator TEXT DEFAULT '?????',
                        Difficulty TEXT DEFAULT '1',
                        AudioFileName TEXT DEFAULT 'unknown.mp3',
                        MD5Hash TEXT DEFAULT '00000000000000000000000000000000',
                        FileName TEXT DEFAULT 'unknown.osu',
                        RankedStatus TEXT DEFAULT Unknown,
                        LastModifiedTime DATETIME DEFAULT '0001-01-01 00:00:00',
                        TotalTime INTEGER DEFAULT 0,
                        AudioPreviewTime INTEGER DEFAULT 0,
                        BeatmapSetId INTEGER DEFAULT -1,
                        Source TEXT DEFAULT '',
                        Tags TEXT DEFAULT '',
                        LastPlayed DATETIME DEFAULT '0001-01-01 00:00:00',
                        FolderName TEXT DEFAULT 'Unknown Folder',
                        UNIQUE (Artist, Title, MD5Hash)
                    );";

                using (var command = new SQLiteCommand(createBeatmaps, connection))
                {
                    command.ExecuteNonQuery();
                }

                string activeSearch = @"CREATE VIRTUAL TABLE IF NOT EXISTS BeatmapFTS USING fts4(
                                          Title, 
                                          Artist, 
                                        );";
                
                using (var command = new SQLiteCommand(activeSearch, connection))
                {
                    command.ExecuteNonQuery();
                }

                string triggerSearchupdate = @"CREATE TRIGGER IF NOT EXISTS Beatmap_Insert_Trigger
                                                    AFTER INSERT ON Beatmap
                                                    BEGIN
                                                        INSERT INTO BeatmapFTS (Title, Artist)
                                                        VALUES (NEW.Title, NEW.Artist);
                                                    END;";

                using (var command = new SQLiteCommand(triggerSearchupdate, connection))
                {
                    command.ExecuteNonQuery();
                }

                string query = @"SELECT COUNT(rowid) as count FROM Beatmap";

                using (var command = new SQLiteCommand(query, connection))
                {
                    using (var reader = command.ExecuteReader())
                    {
                        while (reader.Read())
                        {
                            count = reader.GetInt32(reader.GetOrdinal("count"));
                        }
                    }
                }
            }

            if (count < osuDatabase.BeatmapCount)
            {
                DateTime? LastMapInsert = null;

                using (var connection = new SQLiteConnection(dburl))
                {
                    connection.Open();

                    string query = @"SELECT MAX(LastModifiedTime) as Time FROM Beatmap"; ;
                    using (var command = new SQLiteCommand(query, connection))
                    {
                        using (var reader = command.ExecuteReader())
                        {
                            while (reader.Read())
                            {
                                try
                                {
                                    LastMapInsert = reader.GetDateTime(reader.GetOrdinal("Time"));
                                }

                                catch (Exception e)
                                {
                                    LastMapInsert = null;
                                }
                            }
                        }
                    }

                    int i = 0;
                    int size = osuDatabase.BeatmapCount;
                    Console.CursorVisible = false;
                    foreach (var item in osuDatabase.Beatmaps)
                    {
                        if (LastMapInsert == null || item.LastModifiedTime > LastMapInsert)
                        {
                            insertBeatmap(item);
                            i++;
                            Console.CursorTop -= 1;
                            Console.Write($"Inserted {i}/{size}");
                        }
                    }

                }
            }



        }

        public static List<Song> GetByRecent(int limit, int offset)
        {
            var songs = new List<Song>();

            using (var connection = new SQLiteConnection(dburl))
            {
                connection.Open();

                string query = @"
                     SELECT 
                        MD5Hash, 
                        Title,
                        Artist,
                        TotalTime,
                        Creator,
                        FileName,
                        FolderName,
                        AudioFileName
                    FROM 
                        Beatmap
                    GROUP BY
                        BeatmapSetId
                    ORDER BY
                        LastModifiedTime DESC
                    LIMIT @Limit
                    OFFSET @Offset
                    ";

                using (var command = new SQLiteCommand(query, connection))
                {
                    command.Parameters.AddWithValue("@Limit", limit);
                    command.Parameters.AddWithValue("@Offset", offset);
                    using (var reader = command.ExecuteReader())
                    {
                        while (reader.Read())
                        {
                            string folder = reader.GetString(reader.GetOrdinal("FolderName"));
                            string file = reader.GetString(reader.GetOrdinal("FileName"));
                            string audio = reader.GetString(reader.GetOrdinal("AudioFileName"));

                            var img = Osudb.getBG(folder, file);
                            Song song = new Song(
                                hash: reader.GetString(reader.GetOrdinal("MD5Hash")),
                                name: reader.GetString(reader.GetOrdinal("Title")),
                                artist: reader.GetString(reader.GetOrdinal("Artist")),
                                length: reader.GetInt32(reader.GetOrdinal("TotalTime")),
                            url: $"{folder}/{audio}",
                            previewimage: img,
                                mapper: reader.GetString(reader.GetOrdinal("Creator"))
                            );

                            songs.Add(song);
                        }
                    }
                }
            }
            return songs;
        }

        public void insertBeatmap(OsuParsers.Database.Objects.DbBeatmap beatmap)
        {

            using (var connection = new SQLiteConnection(dburl))
            {

                connection.Open();

                string insertBeatmap = @"
                    INSERT INTO Beatmap (
                        BeatmapId, Artist, ArtistUnicode, Title, TitleUnicode, Creator, Difficulty, 
                        AudioFileName, MD5Hash, FileName, RankedStatus, LastModifiedTime, TotalTime, 
                        AudioPreviewTime, BeatmapSetId, Source, Tags, LastPlayed, FolderName
                    ) VALUES (
                        @BeatmapId, @Artist, @ArtistUnicode, @Title, @TitleUnicode, @Creator, @Difficulty, 
                        @AudioFileName, @MD5Hash, @FileName, @RankedStatus, @LastModifiedTime, @TotalTime, 
                        @AudioPreviewTime, @BeatmapSetId, @Source, @Tags, @LastPlayed, @FolderName
                    );";
                using (var command = new SQLiteCommand(insertBeatmap, connection))
                {
                    command.Parameters.AddWithValue("@BeatmapSetId", beatmap.BeatmapSetId);
                    command.Parameters.AddWithValue("@BeatmapId", beatmap.BeatmapId);
                    command.Parameters.AddWithValue("@Artist", beatmap.Artist);
                    command.Parameters.AddWithValue("@ArtistUnicode", beatmap.ArtistUnicode);
                    command.Parameters.AddWithValue("@Title", beatmap.Title);
                    command.Parameters.AddWithValue("@TitleUnicode", beatmap.TitleUnicode);
                    command.Parameters.AddWithValue("@Creator", beatmap.Creator);
                    command.Parameters.AddWithValue("@Difficulty", beatmap.Difficulty);
                    command.Parameters.AddWithValue("@AudioFileName", beatmap.AudioFileName);
                    command.Parameters.AddWithValue("@MD5Hash", beatmap.MD5Hash);
                    command.Parameters.AddWithValue("@FileName", beatmap.FileName);
                    command.Parameters.AddWithValue("@RankedStatus", beatmap.RankedStatus);
                    command.Parameters.AddWithValue("@LastModifiedTime", beatmap.LastModifiedTime);
                    command.Parameters.AddWithValue("@TotalTime", beatmap.TotalTime);
                    command.Parameters.AddWithValue("@AudioPreviewTime", beatmap.AudioPreviewTime);
                    command.Parameters.AddWithValue("@Source", beatmap.Source);
                    command.Parameters.AddWithValue("@Tags", beatmap.Tags);
                    command.Parameters.AddWithValue("@LastPlayed", beatmap.LastPlayed);
                    command.Parameters.AddWithValue("@FolderName", beatmap.FolderName);

                    int rows = command.ExecuteNonQuery();
                    Console.WriteLine(rows);
                }
            }
        }

        public static Song GetSongByHash(string hash)
        {

            using (var connection = new SQLiteConnection(dburl))
            {
                connection.Open();

                string query = @"
                     SELECT 
                        MD5Hash, 
                        Title,
                        Artist,
                        TotalTime,
                        Creator,
                        FileName,
                        FolderName,
                        AudioFileName
                    FROM Beatmap
                    WHERE MD5Hash = @Hash;
                    ";

                using (var command = new SQLiteCommand(query, connection))
                {
                    command.Parameters.AddWithValue("@Hash", hash);
                    using (var reader = command.ExecuteReader())
                    {
                        while (reader.Read())
                        {

                            string folder = reader.GetString(reader.GetOrdinal("FolderName"));
                            string file = reader.GetString(reader.GetOrdinal("FileName"));
                            string audio = reader.GetString(reader.GetOrdinal("AudioFileName"));

                            var img = Osudb.getBG(folder, file);
                            Song song = new Song(
                                hash: reader.GetString(reader.GetOrdinal("MD5Hash")),
                                name: reader.GetString(reader.GetOrdinal("Title")),
                                artist: reader.GetString(reader.GetOrdinal("Artist")),
                                length: reader.GetInt32(reader.GetOrdinal("TotalTime")),
                            url: $"{folder}/{audio}",
                            previewimage: img,
                                mapper: reader.GetString(reader.GetOrdinal("Creator"))
                            );

                            return song;
                        }
                    }
                }
            }
            return null;
        }

        public static ActiveSearch activeSearch(string query) {
            ActiveSearch search = new ActiveSearch();

            using (var connection = new SQLiteConnection(dburl))
            {
                string q = @"SELECT 
                        MD5Hash, 
                        Title,
                        Artist,
                        TotalTime,
                        Creator,
                        FileName,
                        FolderName,
                        AudioFileName
                    FROM Beatmap
                    WHERE Title LIKE @query
                    OR Artist LIKE @query
                    OR Tags LIKE @query
                    Group By Title
                    LIMIT 15";

                connection.Open();

                using (var command = new SQLiteCommand(q, connection))
                {
                    command.Parameters.AddWithValue("@query", "%" + query + "%");

                    using (var reader = command.ExecuteReader()) {

                        while (reader.Read()) {

                            string folder = reader.GetString(reader.GetOrdinal("FolderName"));
                            string file = reader.GetString(reader.GetOrdinal("FileName"));
                            string audio = reader.GetString(reader.GetOrdinal("AudioFileName"));

                            var img = Osudb.getBG(folder, file);
                            Song song = new Song(
                                hash: reader.GetString(reader.GetOrdinal("MD5Hash")),
                                name: reader.GetString(reader.GetOrdinal("Title")),
                                artist: reader.GetString(reader.GetOrdinal("Artist")),
                                length: reader.GetInt32(reader.GetOrdinal("TotalTime")),
                            url: $"{folder}/{audio}",
                            previewimage: img,
                                mapper: reader.GetString(reader.GetOrdinal("Creator"))
                            );

                            search.Songs.Add(song);
                        }
                    }
                }

                string q2 = @"SELECT 
                        Artist
                    FROM Beatmap
                    WHERE Artist LIKE @query
                    Group By Artist
                    LIMIT 5";

                using (var command = new SQLiteCommand(q2, connection))
                {
                    command.Parameters.AddWithValue("@query", query + "%");

                    using (var reader = command.ExecuteReader())
                    {

                        while (reader.Read())
                        {

                            search.Artist.Add(reader.GetString(reader.GetOrdinal("Artist")));

                        }
                    }
                }
            }

            return search;
        }

        public static List<Song> GetArtistSearch(string query, int limit, int offset) {
            List<Song> songs = new List<Song>();

            query = $"%{query}%";
            using (var connection = new SQLiteConnection(dburl))
            {
                string q = @"SELECT 
                        MD5Hash, 
                        Title,
                        Artist,
                        TotalTime,
                        Creator,
                        FileName,
                        FolderName,
                        AudioFileName
                    FROM Beatmap
                    WHERE Artist LIKE @query
                    Group By Title
                    LIMIT @Limit
                    OFFSET @Offset";

                connection.Open();

                using (var command = new SQLiteCommand(q, connection))
                {
                    command.Parameters.AddWithValue("@query", query);
                    command.Parameters.AddWithValue("@Limit", limit);
                    command.Parameters.AddWithValue("@Offset", offset);

                    using (var reader = command.ExecuteReader())
                    {

                        while (reader.Read())
                        {
                            string folder = reader.GetString(reader.GetOrdinal("FolderName"));
                            string file = reader.GetString(reader.GetOrdinal("FileName"));
                            string audio = reader.GetString(reader.GetOrdinal("AudioFileName"));

                            var img = Osudb.getBG(folder, file);
                            Song song = new Song(
                                hash: reader.GetString(reader.GetOrdinal("MD5Hash")),
                                name: reader.GetString(reader.GetOrdinal("Title")),
                                artist: reader.GetString(reader.GetOrdinal("Artist")),
                                length: reader.GetInt32(reader.GetOrdinal("TotalTime")),
                            url: $"{folder}/{audio}",
                            previewimage: img,
                                mapper: reader.GetString(reader.GetOrdinal("Creator"))
                            );

                            songs.Add(song);
                        }
                    }
                }
            }

                return songs;
        }
    }
}
