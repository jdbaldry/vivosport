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
RETURNING id;

-- name: CreateActivitySession :one
INSERT INTO activity_sessions (
    activity,
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: CreateActivityLap :one
INSERT INTO activity_laps (
    activity,
    message_index,
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: CreateActivityRecord :one
INSERT INTO activity_records (
    activity,
    ts,
    altitude,
    heart_rate,
    cadence,
    distance,
    speed,
    cycles
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: CreateMonitoring :one
INSERT INTO monitorings (
  ts,
  calories,
  cycles,
  distance,
  active_time,
  activity_type,
  activity_sub_type,
  local_ts
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: CreateRecord :one
INSERT INTO records (distance, time)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING id;
