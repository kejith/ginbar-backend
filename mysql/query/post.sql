/* name: GetPosts :many */
SELECT
	* 
FROM
	posts 
WHERE
	deleted_at IS NULL AND
	posts.user_level <= ?
ORDER BY
	posts.id DESC
LIMIT 50;


/* name: GetNewerPosts :many */
SELECT
	* 
FROM
	posts 
WHERE
	deleted_at IS NULL AND
	posts.id >= ? AND
	posts.user_level <= ?
ORDER BY
	posts.id
LIMIT ?;


/* name: GetOlderPosts :many */
SELECT
	* 
FROM
	posts 
WHERE
	deleted_at IS NULL AND
	posts.id <= ? AND
	posts.user_level <= ?
ORDER BY
	posts.id DESC
LIMIT ?;

/* name: GetAllPosts :many */
SELECT
	* 
FROM
	posts;

/* name: GetPost :one */
SELECT
	* 
FROM
	posts 
WHERE
	posts.id = ? AND 
	deleted_at IS NULL AND
	posts.user_level <= ?;

/* name: UpdatePostFiles :exec */
UPDATE
	posts
SET
	filename = ?,
	thumbnail_filename = ?
WHERE
	id = ?;

/* name: GetPostsByUser :many */
SELECT
	* 
FROM
	posts 
WHERE
	user_name = ? AND 
	deleted_at IS NULL AND
	posts.user_level <= ?
ORDER BY posts.id DESC;

/* name: CreatePost :execresult */
INSERT INTO posts 
    (url, filename, thumbnail_filename, user_name, content_type)
VALUES 
    (?, ?, ?, ?, ?);

/* name: DeletePost :exec */
UPDATE posts 
SET deleted_at = NOW() 
WHERE
	id = ?;

/* name: GetVotedPosts :many */
SELECT
	p.*,
	IFNULL(pv.upvoted, 0) as upvoted 
FROM
	posts p
	LEFT JOIN ( SELECT * FROM post_votes WHERE user_id = ? ) AS pv ON pv.post_id = p.id 
WHERE
	p.deleted_at IS NULL AND
	p.user_level <= ?
ORDER BY p.id DESC;

/* name: GetVotedPost :one */
SELECT
	p.*, 
	IFNULL(pv.upvoted, 0) as upvoted 
FROM
	posts p
	LEFT JOIN post_votes AS pv ON pv.post_id = p.id 
	AND pv.user_id = ? AND
	p.user_level <= ?
WHERE
	p.deleted_at IS NULL AND
	p.id = ?;




