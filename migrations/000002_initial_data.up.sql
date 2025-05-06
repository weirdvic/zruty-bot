-- Вставка стандартных настроек
INSERT INTO settings (key, value)
    VALUES 
    ('banAfter', '23')
;

INSERT INTO users (id, username, first_seen_at, check_passed_at)
	VALUES
    (68051500, 'ashenzari', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (482246384, 'zHz', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
;

INSERT INTO admins (user_id)
    VALUES
    (68051500),
    (482246384)
;

INSERT INTO chats (chat_id, title, username, chat_type)
    VALUES
    (-1001287084754, 'NetHack', 'runethack', 'supergroup')
;

INSERT INTO user_chats (user_id, chat_id)
    VALUES
    (68051500, -1001287084754),
    (482246384, -1001287084754)
;
