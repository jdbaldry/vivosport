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

