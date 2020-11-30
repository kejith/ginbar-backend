/* name: CreatePostVote :exec */
INSERT INTO post_votes 
    (user_id, post_id, upvoted)
VALUES 
    (?, ?, ?);

/* name: UpdatePostVote :exec */
UPDATE post_votes
SET upvoted = ?
WHERE id = ?;


/* name: DeletePostVote :exec */
UPDATE post_votes
SET deleted_at = NOW()
WHERE id = ?