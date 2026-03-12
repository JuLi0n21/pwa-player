package main

import (
	v1 "backend/gen"
	"backend/internal/db"
)

func toProtoSongsSqlC(songs []db.Beatmap) []*v1.Song {
	protoSongs := make([]*v1.Song, len(songs))
	for i, s := range songs {
		protoSongs[i] = toProtoSongSqlC(s)
	}

	return protoSongs
}

func toProtoSongSqlC(song db.Beatmap) *v1.Song {
	return &v1.Song{
		BeatmapId: int32(song.Beatmapid.Int64),
		Md5Hash:   song.Md5hash.String,
		Title:     song.Title.String,
		Artist:    song.Artist.String,
		Creator:   song.Creator.String,
		Folder:    song.Folder.String,
		File:      song.File.String,
		Audio:     song.Audio.String,
		TotalTime: song.Totaltime.Int64,
		Image:     extractImageFromFile(osuRoot, song.Folder.String, song.File.String),
	}
}

func toProtoCollectionSqlc(rows []db.GetCollectionByNameRow) *v1.CollectionResponse {
	songs := make([]*v1.Song, len(rows))

	for i, b := range rows {
		songs[i] = toProtoSongSqlC(db.Beatmap{
			Beatmapid: b.Beatmapid,
			Md5hash:   b.Md5hash,
			Title:     b.Title,
			Artist:    b.Artist,
			Creator:   b.Creator,
			Folder:    b.Folder,
			File:      b.File,
			Audio:     b.Audio,
			Totaltime: b.Totaltime,
		})
	}

	return &v1.CollectionResponse{
		Name:  rows[0].Name.String,
		Items: int32(len(rows)),
		Songs: songs,
	}
}

func toProtoCollectionoffsetSqlc(rows []db.GetCollectionByOffsetRow) *v1.CollectionResponse {
	songs := make([]*v1.Song, len(rows))

	for i, b := range rows {
		songs[i] = toProtoSongSqlC(db.Beatmap{
			Beatmapid: b.Beatmapid,
			Md5hash:   b.Md5hash,
			Title:     b.Title,
			Artist:    b.Artist,
			Creator:   b.Creator,
			Folder:    b.Folder,
			File:      b.File,
			Audio:     b.Audio,
			Totaltime: b.Totaltime,
		})
	}

	name := "Empty Collection"
	if len(rows) > 0 {
		name = rows[0].Name.String
	}

	return &v1.CollectionResponse{
		Name:  name,
		Items: int32(len(rows)),
		Songs: songs,
	}
}

func toProtoCollectionsSearchPreview(rows []db.SearchCollectionRow) *v1.SearchCollectionResponse {

	collection := make([]*v1.CollectionPreview, len(rows))
	for i, c := range rows {
		collection[i] = toProtoCollectionPreviewSqlc(c)
	}

	return &v1.SearchCollectionResponse{Collections: collection}
}

func toProtoCollectionPreviewSqlc(row db.SearchCollectionRow) *v1.CollectionPreview {
	return &v1.CollectionPreview{
		Name:  row.Name.String,
		Items: int32(row.Count),
		Image: extractImageFromFile(osuRoot, row.Folder.String, row.File.String),
	}
}
