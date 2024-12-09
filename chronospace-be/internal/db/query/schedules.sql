-- name: CreateSchedule :one
INSERT INTO schedules (
    service_id,
    date,
    time_start,
    time_end,
    status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetScheduleByID :one
SELECT * FROM schedules
WHERE id = $1;

-- name: ListSchedules :many
SELECT * FROM schedules
ORDER BY date, time_start;

-- name: ListSchedulesByService :many
SELECT * FROM schedules
WHERE service_id = $1
ORDER BY date, time_start;

-- name: UpdateSchedule :one
UPDATE schedules
SET 
    service_id = COALESCE($2, service_id),
    date = COALESCE($3, date),
    time_start = COALESCE($4, time_start),
    time_end = COALESCE($5, time_end),
    status = COALESCE($6, status)
WHERE id = $1
RETURNING *;

-- name: DeleteSchedule :exec
DELETE FROM schedules
WHERE id = $1;