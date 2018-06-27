USE blog;

CREATE TABLE `users`
(
  `id`   CHAR(27) NOT NULL,
  `name` VARCHAR(24) NOT NULL,

  PRIMARY KEY (`id`)
);

CREATE TABLE `posts`
(
  `id`         CHAR(27) NOT NULL,
  `user_id`    CHAR(27) NOT NULL,
  `body`       TEXT NOT NULL,
  `created_at` DATETIME NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);
