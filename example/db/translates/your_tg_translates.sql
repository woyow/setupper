INSERT INTO your_tg_translates(key, name)
VALUES
    -- Menu
    ('menu', '{"ru": "🏠 *МЕНЮ*", "en": "🏠 *MENU*"}'),
    ('menu_item_test1', '{"ru": "ТЕСТ КНОПКА 1", "en": "TEST BUTTON 1"}'),
    ('menu_item_test2', '{"ru": "ТЕСТ КНОПКА 2", "en": "TEST BUTTON 2"}'),
    ('menu_item_test3', '{"ru": "ТЕСТ КНОПКА 3", "en": "TEST BUTTON 3"}')

ON CONFLICT (key) DO UPDATE SET name=excluded.name
RETURNING *;
