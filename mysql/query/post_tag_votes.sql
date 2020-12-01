/* name: CreatePostTagVote :exec */
INSERT INTO post_tag_votes 
    (user_id, post_tag_id, upvoted)
VALUES 
    (?, ?, ?);

/* name: UpsertPostTagVote :exec */
REPLACE INTO 
    post_tag_votes(user_id, post_tag_id, upvoted)
VALUE 
    (?, ?, ?);

/* name: DeletePostTagVote :exec */
DELETE FROM post_tag_votes WHERE user_id = ? AND post_tag_id = ?;