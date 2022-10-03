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
