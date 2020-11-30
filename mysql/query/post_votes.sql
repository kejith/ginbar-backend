/* name: CreatePostVote :exec */
INSERT INTO post_votes 
    (user_id, post_id, upvoted)
VALUES 
    (?, ?, ?);

/* name: UpdatePostVote :exec */
UPDATE post_votes
SET upvoted = ?
WHERE id = ?;




/* name: UpsertPostVote :exec */
REPLACE INTO 
    post_votes(user_id, post_id, upvoted)
VALUE 
    (?, ?, ?);

/* name: DeletePostVote :exec */
DELETE FROM post_votes WHERE user_id = ? AND post_id = ?;    