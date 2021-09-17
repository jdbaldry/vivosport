-- [[file:vendor/github.com/tormoder/fit/messages.go::// ActivityMsg represents the activity FIT message type.][Activities]]
CREATE TABLE activities (
    id bigserial UNIQUE,
    start_ts timestamp,
    end_ts timestamp,
    total_timer_time integer,
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
    total_elapsed_time integer,
    total_timer_time integer,
    total_distance integer,
    total_calories smallint,
    avg_speed smallint,
    max_speed smallint,
    avg_heart_rate smallint,
    max_heart_rate smallint,
    PRIMARY KEY (start_ts, end_ts)
);
