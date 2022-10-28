-- +goose Up
CREATE TABLE IF NOT EXISTS `price_history`
(
    `id`         INT NOT NULL AUTO_INCREMENT,
    `best_price` DECIMAL(5, 2),
    `username`   VARCHAR(32),
    `raw_json`   JSON,
    `created_at` DATETIME DEFAULT NOW(),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `price_history`;
