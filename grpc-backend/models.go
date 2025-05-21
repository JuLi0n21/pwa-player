package main

import v1 "backend/gen"

// Song represents a song entity
// @Description Song represents a song with metadata
type Song struct {
	BeatmapID int    `json:"beatmap_id" example:"123456"`
	MD5Hash   string `json:"md5_hash" example:"abcd1234efgh5678"`
	Title     string `json:"title" example:"Shape of You"`
	Artist    string `json:"artist" example:"Ed Sheeran"`
	Creator   string `json:"creator" example:"JohnDoe"`
	Folder    string `json:"folder" example:"osu/Songs/123456"`
	File      string `json:"file" example:"beatmap.osu"`
	Audio     string `json:"audio" example:"audio.mp3"`
	TotalTime int64  `json:"total_time" example:"240"`
	Image     string `json:"image" example:"cover.jpg"`
}

func (song Song) toProto() *v1.Song {
	return &v1.Song{
		BeatmapId: int32(song.BeatmapID),
		Md5Hash:   song.MD5Hash,
		Title:     song.Title,
		Artist:    song.Artist,
		Creator:   song.Creator,
		Folder:    song.Folder,
		File:      song.File,
		Audio:     song.Audio,
		TotalTime: song.TotalTime,
		Image:     song.Image,
	}
}

func toProtoSongs(local []Song) []*v1.Song {
	songs := make([]*v1.Song, len(local))
	for i, s := range local {
		songs[i] = s.toProto()
	}

	return songs
}

// CollectionPreview represents a preview of a song collection
// @Description CollectionPreview contains summary data of a song collection
type CollectionPreview struct {
	Name  string `json:"name" example:"Collection Name"`
	Image string `json:"image" example:"cover.jpg"`
	Items int    `json:"items" example:"10"`
}

func (c CollectionPreview) toProto() *v1.CollectionPreview {
	return &v1.CollectionPreview{
		Name:  c.Name,
		Image: c.Image,
		Items: int32(c.Items),
	}
}

func toProtoCollectionPreview(local []CollectionPreview) []*v1.CollectionPreview {
	collection := make([]*v1.CollectionPreview, len(local))
	for i, c := range local {
		collection[i] = c.toProto()
	}

	return collection
}

// Collection represents a full song collection
// @Description Collection holds a list of songs
type Collection struct {
	Name  string `json:"name" example:"Best of 2023"`
	Items int    `json:"items" example:"15"`
	Songs []Song `json:"songs"`
}

// ActiveSearch represents an active song search query
// @Description ActiveSearch holds search results for a given artist
type ActiveSearch struct {
	Artist string `json:"artist" example:"Ed Sheeran"`
	Songs  []Song `json:"songs"`
}

// Artist represents an active song search query
// @Description Artist holds search results for a given artist
type Artist struct {
	Artist string `json:"artist" example:"Miku"`
	Count  int    `json:"count" example:"21"`
}

func (a Artist) toProto() *v1.Artist {
	return &v1.Artist{
		Artist: a.Artist,
		Items:  int32(a.Count),
	}
}

func toProtoArtist(local []Artist) []*v1.Artist {
	artists := make([]*v1.Artist, len(local))
	for i, a := range local {
		artists[i] = a.toProto()
	}

	return artists
}
