ALTER TABLE users
  ALTER COLUMN best_speed TYPE INT USING best_speed::INT,
  ALTER COLUMN avg_performance TYPE INT USING avg_performance::INT;
