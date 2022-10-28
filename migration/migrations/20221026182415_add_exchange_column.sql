-- +goose Up
ALTER TABLE `price_history` ADD COLUMN `exchange` VARCHAR(24) AFTER `best_price`;

-- +goose Down
ALTER TABLE `price_history` DROP COLUMN `exchange`;
