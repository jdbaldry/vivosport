-- Manual cleanups until I understand how to modify the ACTIVITY/2021-09-20.FIT file correctly.
-- Calories is unchanged but it its only nine minutes of inactivity difference so no biggie.

-- Find the activity.
SELECT id FROM activity WHERE start_ts = timetamp '2021-09-20 19:16:26';

-- Remove the extra records from the end. After 20:37:00, I was driving.
DELETE FROM activity_records WHERE activity = 141 AND ts > timestamp '2021-09-20 20:37:00';

-- Fix all the activity end timestamps.
UPDATE activities SET end_ts = timestamp '2021-09-20 20:37:00' WHERE id = 141;
UPDATE activity_laps SET end_ts = timestamp '2021-09-20 20:37:00' WHERE activity = 141;
UPDATE activity_sessions SET end_ts = timestamp '2021-09-20 20:37:00' WHERE activity = 141;

-- Fix up the laps fields.
UPDATE activity_laps SET avg_heart_rate = (SELECT AVG(heart_rate) FROM activity_records WHERE activity = 141) WHERE activity = 141;
UPDATE activity_laps SET MAX_heart_rate = (SELECT MAX(heart_rate) FROM activity_records WHERE activity = 141) WHERE activity = 141;
UPDATE activity_laps SET total_distance = (SELECT MAX(distance) FROM activity_records WHERE activity = 141) WHERE activity = 141;
UPDATE activity_laps SET avg_speed = (SELECT AVG(speed) FROM activity_records WHERE activity = 141) WHERE activity = 141;
UPDATE activity_laps SET total_timer_time = (SELECT EXTRACT(epoch FROM (end_ts - start_ts)) FROM activity_laps WHERE activity = 141) WHERE activity = 141;
-- Fix up the sessions fields.
UPDATE activity_sessions SET avg_heart_rate = (SELECT AVG(heart_rate) FROM activity_records WHERE activity = 141) WHERE activity = 141;
UPDATE activity_sessions SET MAX_heart_rate = (SELECT MAX(heart_rate) FROM activity_records WHERE activity = 141) WHERE activity = 141;
UPDATE activity_sessions SET total_distance = (SELECT MAX(distance) FROM activity_records WHERE activity = 141) WHERE activity = 141;
UPDATE activity_sessions SET avg_speed = (SELECT AVG(speed) FROM activity_records WHERE activity = 141) WHERE activity = 141;
UPDATE activity_sessions SET total_timer_time = (SELECT EXTRACT(epoch FROM (end_ts - start_ts)) FROM activity_laps WHERE activity = 141) WHERE activity = 141;
