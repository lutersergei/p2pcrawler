-- +goose Up
ALTER TABLE `price_history` ADD COLUMN `exchange` VARCHAR(24) AFTER `max_price`;

-- +goose Down
ALTER TABLE `price_history` DROP COLUMN `exchange`;
