CREATE TABLE comment_votes (
  /* keys */
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  
  /* date */
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME NOT NULL DEFAULT NOW(),
  deleted_at DATETIME DEFAULT NULL,

  /* body */  
  user_id  INT UNSIGNED NOT NULL,
  comment_id INT UNSIGNED NOT NULL,
  upvoted smallint NOT NULL DEFAULT 0
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE comment_votes
  ADD UNIQUE KEY uidx_comment_vote (user_id,comment_id) USING BTREE,
  ADD KEY idx_comment_votes_deleted_at (deleted_at),
  ADD FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  ADD FOREIGN KEY (comment_id) REFERENCES comments(id) ON UPDATE CASCADE ON DELETE CASCADE;


/* 
 *
 * Trigger Updates the Value of Score for the given Comment 
 *
 */

-- CREATE TRIGGER 
-- 	trigger_comment_score_insert_updater
-- AFTER INSERT ON comment_votes 
-- 	FOR EACH ROW
-- 		UPDATE comments 
-- 		SET comments.score = COALESCE(comments.score, 0) + NEW.upvoted
-- 		WHERE
-- 			comments.id = NEW.comment_id;

-- CREATE TRIGGER 
-- 	trigger_comment_score_delete_updater
-- AFTER DELETE ON comment_votes 
-- 	FOR EACH ROW
-- 		UPDATE comments 
-- 		SET comments.score = COALESCE(comments.score, 0) - OLD.upvoted
-- 		WHERE
-- 			comments.id = OLD.comment_id;
			

-- CREATE TRIGGER 
-- 	trigger_comment_score_update_updater
-- AFTER UPDATE ON comment_votes 
-- 	FOR EACH ROW
-- 		UPDATE comments 
-- 		SET comments.score = COALESCE(comments.score, 0) + NEW.upvoted
-- 		WHERE
-- 			comments.id = NEW.comment_id;
			
