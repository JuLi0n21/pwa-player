-- Beatmap Table
CREATE TABLE IF NOT EXISTS Beatmap (
    BeatmapId INTEGER DEFAULT 0,
    Artist TEXT DEFAULT '?????',
    ArtistUnicode TEXT DEFAULT '?????',
    Title TEXT DEFAULT '???????',
    TitleUnicode TEXT DEFAULT '???????',
    Creator TEXT DEFAULT '?????',
    Difficulty TEXT DEFAULT '1',
    Audio TEXT DEFAULT 'unknown.mp3',
    MD5Hash TEXT DEFAULT '00000000000000000000000000000000',
    File TEXT DEFAULT 'unknown.osu',
    RankedStatus TEXT DEFAULT 'Unknown',
    LastModifiedTime INTEGER DEFAULT 0,
    TotalTime INTEGER DEFAULT 0,
    AudioPreviewTime INTEGER DEFAULT 0,
    BeatmapSetId INTEGER DEFAULT -1,
    Source TEXT DEFAULT '',
    Tags TEXT DEFAULT '',
    LastPlayed INTEGER DEFAULT 0,
    Folder TEXT DEFAULT 'Unknown Folder',
    UNIQUE (Artist, Title, MD5Hash, Difficulty, Folder) ON CONFLICT IGNORE
);

CREATE INDEX IF NOT EXISTS idx_beatmap_md5hash ON Beatmap(MD5Hash);
CREATE INDEX IF NOT EXISTS idx_beatmap_lastModifiedTime ON Beatmap(LastModifiedTime);
CREATE INDEX IF NOT EXISTS idx_beatmap_title_artist ON Beatmap(Title, Artist);

-- Collection Table
CREATE TABLE IF NOT EXISTS Collection (
    Name TEXT DEFAULT '',
    MD5Hash TEXT DEFAULT '00000000000000000000000000000000'
);

CREATE INDEX IF NOT EXISTS idx_collection_name ON Collection(Name);
CREATE INDEX IF NOT EXISTS idx_collection_md5hash ON Collection(MD5Hash);
