ALTER TABLE user_money_movement
ADD CONSTRAINT amount_negative CHECK(amount > 0);