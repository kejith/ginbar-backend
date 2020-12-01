CREATE TABLE post_tag_votes (
  /* keys */
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,

  /* body */  
  user_id  INT UNSIGNED NOT NULL,
  post_tag_id INT UNSIGNED NOT NULL,
  upvoted smallint NOT NULL DEFAULT 0
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE post_tag_votes
  ADD UNIQUE KEY uidx_post_tag_vote (user_id, post_tag_id) USING BTREE,
  ADD FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  ADD FOREIGN KEY (post_tag_id) REFERENCES post_tags(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- CREATE TRIGGER 
-- 	trigger_post_tag_score_insert_updater
-- AFTER INSERT ON post_tag_votes 
-- 	FOR EACH ROW
-- 		UPDATE post_tags 
-- 		SET post_tags.score = COALESCE(post_tags.score, 0) + NEW.upvoted
-- 		WHERE
-- 			post_tags.id = NEW.post_tag_id;

-- CREATE TRIGGER 
-- 	trigger_post_tag_score_delete_updater
-- AFTER DELETE ON post_tag_votes 
-- 	FOR EACH ROW
-- 		UPDATE post_tags 
-- 		SET post_tags.score = COALESCE(post_tags.score, 0) - OLD.upvoted
-- 		WHERE
-- 			post_tags.id = OLD.post_tag_id;
			
-- CREATE TRIGGER 
-- 	trigger_post_tag_score_update_updater
-- AFTER UPDATE ON post_tag_votes 
-- 	FOR EACH ROW
-- 		UPDATE post_tags 
-- 		SET post_tags.score = COALESCE(post_tags.score, 0) + NEW.upvoted
-- 		WHERE
-- 			post_tags.id = NEW.post_tag_id;