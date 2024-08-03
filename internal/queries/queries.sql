-- name: CreateUser :exec
INSERT INTO users (
  tgID, login, name, createData, chatID
) VALUES (
  ?1, ?2, ?3, ?4, ?5
) ON CONFLICT (tgID)
DO UPDATE SET login = ?2, name = ?3, chatID = ?5;

-- name: GetUserInfo :one
SELECT login, name, ban 
FROM users
WHERE tgID = ?;

-- name: SetStatus :exec
INSERT INTO chatStatuses (
  tgID, status, annID
) VALUES (
  ?1, ?2, ?3
) ON CONFLICT (tgID)
DO UPDATE SET status = ?2, annID = ?3;

-- name: GetStatus :one
SELECT status 
FROM chatStatuses
WHERE tgID = ?;

-- name: GetAnnId :one
SELECT annID 
FROM chatStatuses
WHERE tgID = ?;

-- name: AddAnnouncement :one
INSERT INTO announcements (
  tgID, txt, chatID
  ) VALUES (
    ?1, ?2, ?3
  )
RETURNING id;

-- name: GetAnnouncement :one
SELECT txt, publicID
FROM announcements
WHERE tgID = ? AND id = ?;

-- name: SetAdminMsgID :exec
UPDATE announcements
SET admMsgId = ?
WHERE id = ?;

-- name: AddPhoto :exec
UPDATE announcements
SET fileID = ?
WHERE id = ?;

-- name: GetAnnouncementOnAdmMsgID :one
SELECT txt, fileID, chatID, id, tgID, publicID
FROM announcements
WHERE admMsgID = ?;

-- name: SetPublicID :exec
UPDATE announcements
SET publicID = ?
WHERE id = ?;

-- name: Ban :exec
UPDATE users
SET ban = true
WHERE tgID = ?;

-- name: UnBan :exec
UPDATE users
SET ban = false
WHERE tgID = ?;
