DELETE FROM settings
WHERE key IN (
    'greetAdminMessage',
    'notAdminMessage',
    'kickMessage',
    'welcomeMessage',
    'challengeMessage'
    );