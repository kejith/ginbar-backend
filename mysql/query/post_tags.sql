/* name: GetTagsByPost :many */
SELECT pt.id, pt.score, pt.post_id, pt.user_id, t.name, IFNULL(ptv.upvoted, 0) upvoted
FROM post_tags pt 
LEFT JOIN tags AS t ON pt.tag_id = t.id
LEFT JOIN post_tag_votes AS ptv ON ptv.post_tag_id = pt.id AND ptv.user_id = ?
WHERE pt.post_id = ?
ORDER BY score DESC;

/* name: GetPostTag :one */
SELECT * FROM post_tags WHERE id = ?;

/* name: AddTagToPost :execresult */
INSERT INTO post_tags (tag_id, post_id, user_id) VALUES (?, ?, ?);

/* name: RemoveTagFromPost :exec */
DELETE FROM post_tags WHERE (tag_id, post_id ) = (?, ?);