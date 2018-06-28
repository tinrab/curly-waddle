USE blog;

CREATE TABLE `users`
(
  `id`   INT         NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,

  PRIMARY KEY (`id`)
);

CREATE TABLE `posts`
(
  `id`         INT      NOT NULL AUTO_INCREMENT,
  `user_id`    INT      NOT NULL,
  `body`       TEXT     NOT NULL,
  `created_at` DATETIME NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
    ON DELETE CASCADE
);
