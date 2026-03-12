-- name: InsertBeatmap :exec
INSERT INTO Beatmap (
    BeatmapId, Artist, ArtistUnicode, Title, TitleUnicode, Creator,
    Difficulty, Audio, MD5Hash, File, RankedStatus, 
    LastModifiedTime, TotalTime, AudioPreviewTime, BeatmapSetId, 
    Source, Tags, LastPlayed, Folder
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetBeatmapByHash :one
SELECT * FROM Beatmap WHERE MD5Hash = ?;

-- name: GetBeatmapCount :one
SELECT COUNT(*) FROM Beatmap;

-- name: GetBeatmapSetCount :one
SELECT COUNT(*) FROM Beatmap GROUP BY BeatmapSetId;

-- name: GetRecentBeatmaps :many
SELECT *
FROM Beatmap GROUP BY Folder ORDER BY LastModifiedTime DESC LIMIT ? OFFSET ?;

-- name: SearchBeatmaps :many
SELECT * 
FROM Beatmap 
WHERE Title LIKE ? OR Artist LIKE ?
GROUP BY Folder
LIMIT ? OFFSET ?;

-- name: SearchArtists :many
SELECT Artist, COUNT(Artist) AS count
FROM Beatmap 
WHERE Artist LIKE ? OR Title LIKE ?
GROUP BY Artist 
LIMIT ? OFFSET ?;
