-- name: CreateUser :one
INSERT INTO users (
  tgID, login, name, createData, chatID
) VALUES (
  $1, $2, $3, $4, $5
) ON CONFLICT (tgID)
DO UPDATE SET login = $2, name = $3, chatID = $5
RETURNING *;

-- name: GetUserInfo :one
SELECT login, name, ban 
FROM users
WHERE tgID = $1;

-- name: SetStatus :one
INSERT INTO chatStatuses (
  tgID, status, annID
) VALUES (
  $1, $2, $3
) ON CONFLICT (tgID)
DO UPDATE SET status = $2, annID = $3
RETURNING *;

-- name: GetStatus :one
SELECT status 
FROM chatStatuses
WHERE tgID = $1;

-- name: GetAnnId :one
SELECT annID 
FROM chatStatuses
WHERE tgID = $1;

-- name: AddAnnouncement :one
INSERT INTO announcements (
  tgID, txt, chatID
  ) VALUES (
    $1, $2, $3
  )
RETURNING id;

-- name: GetAnnouncement :one
SELECT txt, publicID
FROM announcements
WHERE tgID = $1 AND id = $2;

-- name: SetAdminMsgID :exec
UPDATE announcements
SET admMsgId = $1
WHERE id = $2;

-- name: AddPhoto :exec
UPDATE announcements
SET fileID = $2
WHERE id = $1;

-- name: GetAnnouncementOnAdmMsgID :one
SELECT txt, fileID, chatID, id, tgID, publicID
FROM announcements
WHERE admMsgID = $1;

-- name: SetPublicID :exec
UPDATE announcements
SET publicID = $1
WHERE id = $2;

-- name: Ban :exec
UPDATE users
SET ban = true
WHERE tgID = $1;

-- name: UnBan :exec
UPDATE users
SET ban = false
WHERE tgID = $1;
