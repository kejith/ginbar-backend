CREATE TABLE comments (
  /* keys */
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,

  /* date */
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME NOT NULL DEFAULT NOW(),
  deleted_at DATETIME DEFAULT NULL,

  /* body */ 
  content text NOT NULL,
  score int NOT NULL DEFAULT 0,

  /* foreign keys */
  user_name varchar(255) NOT NULL,
  post_id int UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE comments
  ADD KEY idx_comments_deleted_at (deleted_at),
  ADD FOREIGN KEY (user_name) REFERENCES users(name) ON UPDATE CASCADE ON DELETE CASCADE,
  ADD FOREIGN KEY (post_id) REFERENCES posts(id) ON UPDATE CASCADE ON DELETE CASCADE;