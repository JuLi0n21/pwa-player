-- name: InsertCollection :exec
INSERT INTO Collection (Name, MD5Hash) VALUES (?, ?);

-- name: GetCollectionCountByName :one
SELECT COUNT(*) FROM Collection WHERE Name = ?;

-- name: GetCollections :many
SELECT c.Name, COUNT(b.MD5Hash) AS Count, MIN(b.Folder) AS Folder, MIN(b.File) AS File
FROM Collection c
JOIN Beatmap b ON c.MD5Hash = b.MD5Hash
WHERE c.Name LIKE ?
GROUP BY c.Name
LIMIT ? OFFSET ?;

-- name: GetCollection :many
SELECT 
  c.Name, 
  b.BeatmapId, 
  b.MD5Hash, 
  b.Title, 
  b.Artist, 
  b.Creator, 
  b.Folder, 
  b.File, 
  b.Audio, 
  b.TotalTime
FROM Collection c
JOIN Beatmap b ON c.MD5Hash = b.MD5Hash
JOIN (
  SELECT DISTINCT Name
  FROM Collection
  ORDER BY Name
  LIMIT 1 OFFSET ?
) selected_name ON c.Name = selected_name.Name
LIMIT ? OFFSET ?;

-- name: GetCollectionByName :many
SELECT c.Name, b.BeatmapId, b.MD5Hash, b.Title, b.Artist, b.Creator, b.Folder, b.File, b.Audio, b.TotalTime
FROM Collection c
JOIN Beatmap b ON c.MD5Hash = b.MD5Hash
WHERE c.Name = ?
LIMIT ? OFFSET ?;
