CREATE TABLE IF NOT EXISTS your_tg_translates (
    key varchar PRIMARY KEY,         -- 'hello'
    name jsonb NOT NULL DEFAULT '{}' -- {"ru": "Привет", "en": "Hello"}
);

CREATE TABLE IF NOT EXISTS your_tg_users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id varchar NOT NULL,
    UNIQUE (chat_id)
);

CREATE TABLE IF NOT EXISTS your_tg_chat_states (
    chat_id varchar PRIMARY KEY,
    state varchar NOT NULL DEFAULT 'default'
);

CREATE TABLE IF NOT EXISTS your_tg_user_bans (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tg_id uuid NOT NULL,
    reason varchar NOT NULL DEFAULT '',
    additional_info varchar NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    end_at timestamptz NOT NULL,
    CONSTRAINT fk_your_tg_user FOREIGN KEY (tg_id) REFERENCES your_tg_users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS your_tg_user_bans_search_idx ON your_tg_user_bans(tg_id, end_at);

CREATE OR REPLACE FUNCTION your_tg_check_ban(req_your_tg_chat varchar) RETURNS bool
    LANGUAGE plpgsql AS
$$
DECLARE
    is_banned_out bool;
BEGIN
    SELECT EXISTS (
        SELECT 1
        FROM your_tg_user_bans ytub
        WHERE ytub.tg_id=(select id from your_tg_users where chat_id=req_your_tg_chat)
          AND ytub.end_at >= current_timestamp
    ) INTO is_banned_out;

    RETURN is_banned_out;
END
$$;

-- translates for test
INSERT INTO your_tg_translates(key, name)
VALUES
    -- Menu
    ('menu', '{"ru": "🏠 *МЕНЮ*", "en": "🏠 *MENU*"}'),
    ('menu_item_test1', '{"ru": "ТЕСТ КНОПКА 1", "en": "TEST BUTTON 1"}'),
    ('menu_item_test2', '{"ru": "ТЕСТ КНОПКА 2", "en": "TEST BUTTON 2"}'),
    ('menu_item_test3', '{"ru": "ТЕСТ КНОПКА 3", "en": "TEST BUTTON 3"}')

ON CONFLICT (key) DO UPDATE SET name=excluded.name
RETURNING *;
