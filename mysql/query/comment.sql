/* name: GetComments :many */
SELECT *
FROM comments
WHERE deleted_at is NULL
ORDER BY id;

/* name: GetComment :one */
SELECT *
FROM comments
WHERE id = ? AND deleted_at is NULL;

/* name: GetCommentsByPost :many */
SELECT *
FROM comments
WHERE post_id = ? AND deleted_at is NULL;

/* name: CreateComment :exec */
INSERT INTO comments 
    (content, user_name, post_id)
VALUES 
    (?, ?, ?);

/* name: GetLatestComment :one */
SELECT * 
FROM comments
WHERE user_name = ?
ORDER BY created_at DESC
LIMIT 1;

/* name: DeleteComment :exec */
UPDATE comments
SET deleted_at = NOW()
WHERE id = ?;

/* name: GetVotedComments :many */
SELECT
	c.*,
	IFNULL(cv.upvoted, 0) upvoted
FROM
	comments c
	LEFT JOIN comment_votes AS cv ON cv.comment_id = c.id 
	AND cv.user_id = ?
WHERE
	c.deleted_at IS NULL AND
	c.post_id = ?;

/* name: GetVotedComment :one */
SELECT
	c.*,
	cv.upvoted
FROM
	comments c
	LEFT JOIN comment_votes AS cv ON cv.comment_id = c.id 
	AND cv.user_id = ?
WHERE
	c.deleted_at IS NULL AND
	c.id = ?;