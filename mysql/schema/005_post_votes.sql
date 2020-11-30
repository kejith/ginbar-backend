CREATE TABLE post_votes (
  /* keys */
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,

  /* date */
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME NOT NULL DEFAULT NOW(),
  deleted_at DATETIME DEFAULT NULL,

  /* body */
  upvoted SMALLINT NOT NULL DEFAULT 0,

  /* foreign keys */
  user_id INT UNSIGNED NOT NULL,
  post_id INT UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE post_votes
  ADD UNIQUE KEY uidx_post_vote (user_id, post_id),
  ADD KEY idx_post_votes_deleted_at (deleted_at),
  ADD FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  ADD FOREIGN KEY (post_id) REFERENCES posts(id) ON UPDATE CASCADE ON DELETE RESTRICT;


/* 
 *
 * Trigger Updates the Value of Score for the given Posts 
 *
 */

-- CREATE TRIGGER 
-- 	trigger_post_score_insert_updater
-- AFTER INSERT ON post_votes 
-- 	FOR EACH ROW
-- 		UPDATE posts 
-- 		SET posts.score = COALESCE(posts.score, 0) + NEW.upvoted
-- 		WHERE
-- 			posts.id = NEW.post_id;

-- CREATE TRIGGER 
-- 	trigger_post_score_delete_updater
-- AFTER DELETE ON post_votes 
-- 	FOR EACH ROW
-- 		UPDATE posts 
-- 		SET posts.score = COALESCE(posts.score, 0) - OLD.upvoted
-- 		WHERE
-- 			posts.id = OLD.post_id;
			
-- CREATE TRIGGER 
-- 	trigger_post_score_update_updater
-- AFTER UPDATE ON post_votes 
-- 	FOR EACH ROW
-- 		UPDATE posts 
-- 		SET posts.score = COALESCE(posts.score, 0) + NEW.upvoted
-- 		WHERE
-- 			posts.id = NEW.post_id;