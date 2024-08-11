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