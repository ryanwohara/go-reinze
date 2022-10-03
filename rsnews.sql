CREATE TABLE `reinze`.`rsnews` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NOT NULL,
  `url` VARCHAR(125) NOT NULL,
  `hash_id` VARCHAR(5) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT now(),
  `runescape` ENUM("oldschool", "runescape3") NOT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `reinze`.`news` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NOT NULL,
  `url` VARCHAR(255) NOT NULL,
  `hash_id` VARCHAR(5) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`));

CREATE TABLE `reinze`.`to_post` (
  `id` int NOT NULL AUTO_INCREMENT,
  `path_on_disk` varchar(125) NOT NULL,
  `scheduled_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `posted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `reinze`.`twitter` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `text` VARCHAR(255) NOT NULL,
  `tweet_id` VARCHAR(20) NOT NULL,
  PRIMARY KEY (`id`));