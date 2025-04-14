-- Удаление стандартных настроек
DELETE FROM settings WHERE key IN ('ban_after', 'groups');

-- Удаление пользователей
DELETE FROM admins;
DELETE FROM users;
DELETE FROM chats;