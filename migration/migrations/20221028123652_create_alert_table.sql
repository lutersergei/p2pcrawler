-- +goose Up
CREATE TABLE IF NOT EXISTS `alert_request`
(
    `id`         INT NOT NULL AUTO_INCREMENT,
    `price`      DECIMAL(10, 2),
    `username`   VARCHAR(32),
    `exchange`   VARCHAR(32),
    `move_type`  VARCHAR(4),
    `deal_type`  VARCHAR(4),
    `status`     VARCHAR(6),
    `created_at` DATETIME DEFAULT NOW(),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `alert_request`;
