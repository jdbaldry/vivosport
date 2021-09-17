-- name: CreateActivity :one
INSERT INTO activities (
    start_ts,
    end_ts,
    total_timer_time,
    num_sessions,
    type,
    event,
    event_type,
    local_ts,
    event_group
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT DO NOTHING
RETURNING *;
