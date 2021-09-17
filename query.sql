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

-- name: CreateSession :one
INSERT INTO sessions (
    start_ts,
    end_ts,
    event,
    event_type,
    sport,
    sub_sport,
    total_elapsed_time,
    total_timer_time,
    total_distance,
    total_calories,
    avg_speed,
    max_speed,
    avg_heart_rate,
    max_heart_rate
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
ON CONFLICT DO NOTHING
RETURNING *;
