-- name: CreateUser :exec
INSERT INTO users (
  tgID, login, name, createData, chatID
) VALUES (
  ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE login = VALUES(login), name = VALUES(name), chatID = VALUES(chatID);

-- name: GetUserInfo :one
SELECT login, name, ban 
FROM users
WHERE tgID = ?;

-- name: SetStatus :exec
INSERT INTO chatStatuses (
  tgID, status, annID
) VALUES (
  ?, ?, ?
) ON DUPLICATE KEY UPDATE status = VALUES(status), annID = VALUES(annID);

-- name: GetStatus :one
SELECT status 
FROM chatStatuses
WHERE tgID = ?;

-- name: GetAnnId :one
SELECT annID 
FROM chatStatuses
WHERE tgID = ?;

-- name: AddAnnouncement :execresult
INSERT INTO announcements (
  tgID, txt, chatID
  ) VALUES (
    ?, ?, ?
  );

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
