package main

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

// CollectionPreview represents a preview of a song collection
// @Description CollectionPreview contains summary data of a song collection
type CollectionPreview struct {
	Name  string `json:"name" example:"Collection Name"`
	Image string `json:"image" example:"cover.jpg"`
	Items int    `json:"items" example:"10"`
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
