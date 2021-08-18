CREATE TABLE `reinze`.`to_post` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `path_on_disk` VARCHAR(125) NOT NULL,
  `scheduled_at` TIMESTAMP NOT NULL DEFAULT now(),
  `created_at` TIMESTAMP NOT NULL DEFAULT now(),
  `posted_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`));
