ALTER TABLE user_money_movement
ADD CONSTRAINT date_not_null CHECK(date IS NOT NULL);