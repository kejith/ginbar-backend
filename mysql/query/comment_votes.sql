/* name: CreateCommentVote :exec */
INSERT INTO comment_votes 
    (user_id, comment_id, upvoted)
VALUES 
    (?, ?, ?);

/* name: UpdateCommentVote :exec */
UPDATE comment_votes
SET upvoted = ?
WHERE id = ?;

/* name: UpsertCommentVote :exec */
REPLACE INTO 
    comment_votes(user_id, comment_id, upvoted)
VALUE 
    (?, ?, ?);

/* name: DeleteCommentVote :exec */
DELETE FROM comment_votes WHERE user_id = ? AND comment_id = ?;