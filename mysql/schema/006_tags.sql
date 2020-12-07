CREATE TABLE tags (
  /* keys */
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_level INT UNSIGNED NOT NULL DEFAULT 0,

  /* body */  
  name VARCHAR(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE tags
  ADD UNIQUE KEY uidx_tags_name (name),
  ADD KEY idx_tags_userlevel (user_level);