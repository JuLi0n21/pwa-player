-- name: InsertCollection :exec
INSERT INTO Collection (Name, MD5Hash) VALUES (?, ?);

-- name: GetCollectionCountByName :one
SELECT COUNT(*) FROM Collection WHERE Name = ?;

-- name: SearchCollection :many
SELECT 
  c.Name, 
  COUNT(b.MD5Hash) AS Count, 
  b.Folder, 
  b.File
  FROM Collection c
JOIN Beatmap b ON c.MD5Hash = b.MD5Hash
WHERE c.Name LIKE ?
GROUP BY c.Name
LIMIT ? OFFSET ?;

-- name: GetCollectionByOffset :many
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
WHERE c.Name = (
  SELECT Name
  FROM Collection
  GROUP BY Name
  ORDER BY Name
  LIMIT 1 OFFSET @index
)
LIMIT @limit OFFSET @offset;

-- name: GetCollectionByName :many
SELECT c.Name, b.BeatmapId, b.MD5Hash, b.Title, b.Artist, b.Creator, b.Folder, b.File, b.Audio, b.TotalTime
FROM Collection c
JOIN Beatmap b ON c.MD5Hash = b.MD5Hash
WHERE c.Name = ?
LIMIT ? OFFSET ?;
