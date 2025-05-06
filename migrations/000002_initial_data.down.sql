-- Удаление стандартных настроек
DELETE FROM settings WHERE key IN ('banAfter');

-- Удаление пользователей
DELETE FROM admins;
DELETE FROM users;
DELETE FROM chats;