-- name: GetRTFromUserID :one
SELECT token FROM refresh_tokens WHERE user_id = $1;