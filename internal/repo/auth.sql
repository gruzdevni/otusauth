-- name: IsAuth :one
SELECT * FROM logged_in WHERE user_guid = @guid AND expiry > now();

-- name: InsertSession :exec
INSERT INTO logged_in (user_guid, expiry)
VALUES (@user_guid, now() + INTERVAL '48 hours') ON CONFLICT (user_guid) DO
UPDATE
SET expiry = now() + INTERVAL '48 hours';
