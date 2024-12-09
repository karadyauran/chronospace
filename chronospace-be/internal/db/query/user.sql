-- Insert a new user
INSERT INTO users (username, email, password) 
VALUES ($1, $2, $3)
RETURNING *;

-- Select a user by id
SELECT * FROM users WHERE id = $1;

-- Update a user by id
UPDATE users 
SET username = $2, email = $3, password = $4 
WHERE id = $1
RETURNING *;

-- Delete a user by id
DELETE FROM users WHERE id = $1
RETURNING *;