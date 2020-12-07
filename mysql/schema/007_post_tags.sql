CREATE TABLE post_tags (
  /* keys */
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,

  /* body */
  score SMALLINT NOT NULL DEFAULT 0,

  /* foreign keys */
  tag_id INT UNSIGNED NOT NULL,
  post_id INT UNSIGNED NOT NULL,
  user_id INT UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE post_tags
  ADD UNIQUE KEY uidx_post_tags (tag_id, post_id),
  ADD FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  ADD FOREIGN KEY (tag_id) REFERENCES tags(id) ON UPDATE CASCADE ON DELETE CASCADE,
  ADD FOREIGN KEY (post_id) REFERENCES posts(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- CREATE TRIGGER 
-- 	trigger_update_userlevel_on_create_post_tag
-- AFTER INSERT ON post_tags 
-- 	FOR EACH ROW
-- 		UPDATE posts p
-- 	  LEFT JOIN post_tags ptags ON p.id = NEW.post_id
-- 		LEFT JOIN tags t ON t.id = NEW.tag_id
-- 		SET p.user_level = GREATEST(t.user_level, p.user_level)
-- 		WHERE
-- 			p.id = NEW.post_id;
