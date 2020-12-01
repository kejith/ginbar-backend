/* name: GetTagsByPost :many */
SELECT pt.id, pt.score, t.name FROM post_tags pt 
LEFT JOIN tags AS t ON pt.tag_id = t.id 
WHERE pt.post_id = ?
ORDER BY score DESC;

/* name: AddTagToPost :exec */
INSERT INTO post_tags (tag_id, post_id) VALUES (?, ?);

/* name: RemoveTagFromPost :exec */
DELETE FROM post_tags WHERE (tag_id, post_id ) = (?, ?);