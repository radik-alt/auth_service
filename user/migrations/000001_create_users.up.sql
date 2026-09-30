CREATE TABLE users (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email          TEXT        NOT NULL,
    username       TEXT        NOT NULL,
    phone          TEXT,
    avatar_url     TEXT        NOT NULL DEFAULT '',
    bio            TEXT        NOT NULL DEFAULT '',
    password_hash  TEXT        NOT NULL,
    status         TEXT        NOT NULL DEFAULT 'active'
                               CHECK (status IN ('active', 'blocked', 'deleted')),
    email_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ,

    CONSTRAINT users_email_lowercase CHECK (email = lower(email))
);

-- Уникальность только среди неудалённых: после мягкого удаления email и username снова свободны
CREATE UNIQUE INDEX users_email_uniq    ON users (email)           WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX users_username_uniq ON users (lower(username)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX users_phone_uniq    ON users (phone)           WHERE deleted_at IS NULL AND phone IS NOT NULL;

-- Курсорная пагинация ListUsers
CREATE INDEX users_created_at_id_idx ON users (created_at, id);
