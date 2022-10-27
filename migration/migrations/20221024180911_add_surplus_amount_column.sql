-- +goose Up
ALTER TABLE `price_history` ADD COLUMN `surplus_amount` DECIMAL(12,2) AFTER `max_price`;

-- +goose Down
ALTER TABLE `price_history` DROP COLUMN `surplus_amount`;
