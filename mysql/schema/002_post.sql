CREATE TABLE posts (
  /* keys */
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,

  /* date */
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME NOT NULL DEFAULT NOW(),
  deleted_at DATETIME DEFAULT NULL,

  /* body */  
  url TEXT NOT NULL,
  filename VARCHAR(255) NOT NULL,
  thumbnail_filename VARCHAR(255) NOT NULL,
  content_type VARCHAR(255) NOT NULL,
  score int NOT NULL DEFAULT 0,

  /* foreign key*/
  user_name VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE posts
  ADD FOREIGN KEY (user_name) REFERENCES users(name) ON UPDATE CASCADE ON DELETE RESTRICT;
