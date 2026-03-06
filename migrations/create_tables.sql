CREATE TABLE account (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username text NOT NULL CHECK (length(username) >= 4 AND length(username) <= 32 AND username NOT LIKE '% %'),
    token_version integer NOT NULL DEFAULT 1,
    password text NOT NULL CHECK (length(password) >= 8 AND length(password) <= 40), 
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,

    CONSTRAINT user_username_unique UNIQUE (username)
);

CREATE TABLE note (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    title text NOT NULL CHECK (length(trim(title)) > 0 AND length(title) < 40),
    parent_id uuid,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,

    CONSTRAINT user_id_fk FOREIGN KEY (user_id)
        REFERENCES account (id) ON DELETE CASCADE,
    CONSTRAINT note_parent_id_fkey FOREIGN KEY (parent_id)
        REFERENCES note (id) ON DELETE CASCADE,
    CONSTRAINT note_user_title_parent_unique UNIQUE (user_id, title, parent_id),
    CONSTRAINT note_no_self_parent CHECK (id != parent_id)
);

CREATE TABLE block_type (
    id smallint PRIMARY KEY CHECK (id > 0),
    name text NOT NULL CHECK (length(name) > 0 AND length(name) <= 50),

    CONSTRAINT block_type_name_unique UNIQUE (name)
);

CREATE TABLE block (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    note_id uuid NOT NULL,
    position integer NOT NULL CHECK (position >= 0),
    block_type_id smallint NOT NULL,
    content text,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,

    CONSTRAINT block_note_id_fk FOREIGN KEY (note_id)
        REFERENCES note (id) ON DELETE CASCADE,
    CONSTRAINT block_block_type_id_fk FOREIGN KEY (block_type_id)
        REFERENCES block_type (id) ON DELETE RESTRICT,
    CONSTRAINT block_note_position_unique UNIQUE (note_id, position)
);

CREATE TABLE block_state (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    block_id uuid NOT NULL,
    formatting text NOT NULL CHECK (length(formatting) > 0),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,

    CONSTRAINT block_state_block_id_fk FOREIGN KEY (block_id)
        REFERENCES block (id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION set_timestamps_current_time()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        NEW.created_at = CURRENT_TIMESTAMP;
        NEW.updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_OP = 'UPDATE' THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_account_timestamps 
    BEFORE INSERT OR UPDATE ON account
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps_current_time();

CREATE TRIGGER set_note_timestamps 
    BEFORE INSERT OR UPDATE ON note
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps_current_time();

CREATE TRIGGER set_block_timestamps 
    BEFORE INSERT OR UPDATE ON block
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps_current_time();

CREATE TRIGGER set_block_timestamps 
    BEFORE INSERT OR UPDATE ON block_state
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps_current_time();
