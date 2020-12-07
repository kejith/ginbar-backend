CREATE TABLE users (
  /* keys */
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  
  /* date */
  created_at DATETIME NOT NULL DEFAULT NOW(),
  updated_at DATETIME NOT NULL DEFAULT NOW(),
  deleted_at DATETIME DEFAULT NULL,

  /* body */  
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  level INT UNSIGNED NOT NULL DEFAULT 1

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/* KEYS */
ALTER TABLE users
  ADD UNIQUE KEY name (name),
  ADD UNIQUE KEY email (email),
  ADD KEY idx_user_level (level),
  ADD KEY idx_users_deleted_at (deleted_at);



  