using OsuParsers.Database.Objects;
using OsuParsers.Enums.Database;
using OsuParsers.Enums;

public class Song{
    public string hash {get; set;}
    public string name {get; set;}
    public string artist {get; set;}
    public int length {get; set;}
    public string url { get; set; }
    public string previewimage {get; set;}
    public string mapper {get; set;}

    public Song(string hash, string name, string artist, int length, string url, string previewimage, string mapper) { 
        this.hash = hash; this.name = name; this.artist = artist; this.length = length; this.url = url; this.previewimage = previewimage; this.mapper = mapper;
    }

}

public class CollectionPreview{
    public int index { get; set;}
    public string name {get; set;}
    public int length {get; set;}
    public string previewimage {get; set;}

    private CollectionPreview() { }

    public CollectionPreview(int index, string name, string previewimage, int length) {
        this.index = index; this.name = name; this.previewimage = previewimage; this.length = length;
    }
    
}
public class Collection{
    public string name {get; set;}
    public int length {get; set;}
    public List<Song> songs { get; set;}

    private Collection() { }

    public Collection(string name, int length, List<Song> songs) { 
        this.name = name; this.length = length; this.songs = songs;
    }
    
}

public class ActiveSearch{
    public List<string> Artist { get; set; } = new List<string>();
    public List<Song> Songs { get; set; } = new List<Song>();
}

public class Beatmap
{
    public string Artist { get; set; }
    public string ArtistUnicode { get; set; }
    public string Title { get; set; }
    public string TitleUnicode { get; set; }
    public string Creator { get; set; }
    public string Difficulty { get; set; }
    public string AudioFileName { get; set; }
    public string MD5Hash { get; set; }
    public string FileName { get; set; }
    public RankedStatus RankedStatus { get; set; }
    public DateTime LastModifiedTime { get; set; }
    public int TotalTime { get; set; }
    public int AudioPreviewTime { get; set; }
    public int BeatmapId { get; set; }
    public int BeatmapSetId { get; set; }
    public string Source { get; set; }
    public string Tags { get; set; }
    public DateTime LastPlayed { get; set; }
    public string FolderName { get; set; }
}

public class BeatmapSet {
    public int BeatmapSetId { get; set; }
    public string FolderName { get; set; }
    public string Creator { get; set; }
    public DateTime LastModifiedTime { get; set; }
    public List<Beatmap> Beatmaps { get; private set; } = new List<Beatmap>();

    public void AddBeatmap(Beatmap beatmap)
    {
        beatmap.BeatmapSetId = this.BeatmapSetId;
        Beatmaps.Add(beatmap);
    }
}