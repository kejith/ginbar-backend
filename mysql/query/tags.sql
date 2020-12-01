/* name: GetTags :many */
SELECT * FROM tags ORDER BY	id;

/* name: GetTag :one */
SELECT * FROM tags WHERE tags.id = ? LIMIT 1;

/* name: GetTagByName :one */
SELECT * FROM tags WHERE tags.name = ? LIMIT 1;

/* name: CreateTag :exec */
INSERT INTO tags (name) VALUES (?);

/* name: DeleteTag :exec */
DELETE FROM tags WHERE tags.id = ?;   

/* name: DeleteTagByName :exec */
DELETE FROM tags WHERE tags.name = ?;