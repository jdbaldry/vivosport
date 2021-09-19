-- [[file:vendor/github.com/tormoder/fit/messages.go::// ActivityMsg represents the activity FIT message type.][Activities]]
CREATE TABLE activities (
    id bigserial UNIQUE,
    start_ts timestamp,
    end_ts timestamp,
    total_timer_time double precision, -- double precision is used because the scaled values returns a float64
    num_sessions integer,
    type integer,
    event smallint,
    event_type smallint,
    local_ts timestamp,
    event_group smallint,
    PRIMARY KEY (start_ts, end_ts)
);

-- [[file:vendor/github.com/tormoder/fit/messages.go::// SessionMsg represents the session FIT message type.][Sessions]]
CREATE TABLE sessions (
    id bigserial UNIQUE,
    start_ts timestamp,
    end_ts timestamp,
    event smallint,
    event_type smallint,
    sport smallint,
    sub_sport smallint,
    total_elapsed_time double precision, -- double precision is used because the scaled values returns a float64
    total_timer_time double precision, -- double precision is used because the scaled values returns a float64
    total_distance double precision, -- double precision is used because the scaled values returns a float64
    total_calories smallint,
    avg_speed double precision, -- double precision is used because the scaled values returns a float64
    max_speed double precision, -- double precision is used because the scaled values returns a float64
    avg_heart_rate smallint,
    max_heart_rate smallint,
    PRIMARY KEY (start_ts, end_ts)
);

-- [[file:vendor/github.com/tormoder/fit/messages.go::// MonitoringMsg represents the monitoring FIT message type.][Monitorings]]
CREATE TABLE monitorings (
  id bigserial UNIQUE,
  ts timestamp,
  cycles integer,
  calories smallint,
  distance double precision, -- double precision is used because the scaled values return a float64
  active_time double precision,  -- double precision is used because the scaled values return a float64
  activity_type smallint,
  activity_sub_type smallint,
  local_ts timestamp,
  PRIMARY KEY (ts, activity_type, activity_sub_type)
);
