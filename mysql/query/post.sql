/* name: GetPosts :many */
SELECT
	* 
FROM
	posts 
WHERE
	deleted_at IS NULL 
ORDER BY
	id;

/* name: GetPost :one */
SELECT
	* 
FROM
	posts 
WHERE
	posts.id = ? 
	AND deleted_at IS NULL;

/* name: GetPostsByUser :many */
SELECT
	* 
FROM
	posts 
WHERE
	user_name = ? 
	AND deleted_at IS NULL;

/* name: CreatePost :exec */
INSERT INTO posts 
    (url, image, user_name)
VALUES 
    (?, ?, ?);

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
	p.deleted_at IS NULL;

/* name: GetVotedPost :one */
SELECT
	p.*, 
	IFNULL(pv.upvoted, 0) as upvoted 
FROM
	posts p
	LEFT JOIN post_votes AS pv ON pv.post_id = p.id 
	AND pv.user_id = ? 
WHERE
	p.deleted_at IS NULL AND
	p.id = ?;