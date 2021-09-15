-- [[file:vendor/github.com/tormoder/fit/messages.go::// ActivityMsg represents the activity FIT message type.][Activities]]
CREATE TABLE activities (
    id bigserial PRIMARY KEY,
    ts timestamp,
    total_timer_time integer,
    num_sessions integer,
    type integer,
    event integer,
    event_type integer,
    local_ts timestamp,
    event_group integer
);

