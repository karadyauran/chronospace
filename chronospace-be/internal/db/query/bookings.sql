-- name: CreateBooking :one
INSERT INTO bookings (
    user_id,
    service_id,
    date,
    time,
    status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBooking :one
SELECT * FROM bookings
WHERE id = $1;

-- name: ListBookings :many
SELECT * FROM bookings
ORDER BY date, time;

-- name: ListBookingsByUser :many
SELECT * FROM bookings
WHERE user_id = $1
ORDER BY date, time;

-- name: UpdateBooking :one
UPDATE bookings
SET 
    date = $2,
    time = $3,
    status = $4
WHERE id = $1
RETURNING *;

-- name: DeleteBooking :exec
DELETE FROM bookings
WHERE id = $1;