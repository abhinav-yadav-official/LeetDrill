-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id                          BIGSERIAL PRIMARY KEY,
    email                       TEXT UNIQUE,
    password_hash               TEXT NOT NULL DEFAULT '',
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    leetcode_username           TEXT,
    leetcode_session_encrypted  BYTEA,
    leetcode_csrf_encrypted     BYTEA,
    cookie_updated_at           TIMESTAMPTZ,
    cookie_valid                BOOLEAN NOT NULL DEFAULT FALSE,
    vacation_until              TIMESTAMPTZ
);

CREATE TABLE problems (
    id                      BIGSERIAL PRIMARY KEY,
    leetcode_slug           TEXT UNIQUE NOT NULL,
    leetcode_question_id    TEXT,
    leetcode_frontend_id    TEXT,
    title                   TEXT NOT NULL,
    difficulty              TEXT NOT NULL CHECK (difficulty IN ('Easy','Medium','Hard')),
    url                     TEXT NOT NULL,
    content_html            TEXT,
    topic_tags              JSONB NOT NULL DEFAULT '[]'::jsonb,
    ac_rate                 DOUBLE PRECISION,
    paid_only               BOOLEAN NOT NULL DEFAULT FALSE,
    synced_at               TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX problems_topic_tags_gin ON problems USING GIN (topic_tags);
CREATE INDEX problems_difficulty_idx ON problems (difficulty);

CREATE TABLE patterns (
    id          BIGSERIAL PRIMARY KEY,
    slug        TEXT UNIQUE NOT NULL,
    name        TEXT NOT NULL,
    description TEXT
);

CREATE TABLE problem_patterns (
    problem_id BIGINT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    pattern_id BIGINT NOT NULL REFERENCES patterns(id) ON DELETE CASCADE,
    PRIMARY KEY (problem_id, pattern_id)
);

CREATE INDEX problem_patterns_pattern_idx ON problem_patterns (pattern_id);

CREATE TABLE user_problems (
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    problem_id          BIGINT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    ease_factor         NUMERIC(4,2) NOT NULL DEFAULT 2.5,
    interval_days       INT NOT NULL DEFAULT 0,
    next_due_at         TIMESTAMPTZ,
    last_attempted_at   TIMESTAMPTZ,
    total_attempts      INT NOT NULL DEFAULT 0,
    clean_solves        INT NOT NULL DEFAULT 0,
    total_fails         INT NOT NULL DEFAULT 0,
    streak              INT NOT NULL DEFAULT 0,
    status              TEXT NOT NULL DEFAULT 'new'
                        CHECK (status IN ('new','learning','review','mastered','leech')),
    PRIMARY KEY (user_id, problem_id)
);

CREATE INDEX user_problems_due_idx
    ON user_problems (user_id, next_due_at)
    WHERE status NOT IN ('leech','new');

CREATE INDEX user_problems_status_idx ON user_problems (user_id, status);

CREATE TABLE attempts (
    id                          BIGSERIAL PRIMARY KEY,
    user_id                     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    problem_id                  BIGINT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    started_at                  TIMESTAMPTZ,
    completed_at                TIMESTAMPTZ NOT NULL DEFAULT now(),
    verdict                     TEXT NOT NULL
                                CHECK (verdict IN ('AC','WA','TLE','MLE','RE','CE')),
    submission_count_in_session INT NOT NULL DEFAULT 1,
    time_taken_sec              INT NOT NULL DEFAULT 0,
    runtime_ms                  INT,
    memory_kb                   INT,
    language                    TEXT,
    code                        TEXT,
    derived_rating              TEXT NOT NULL
                                CHECK (derived_rating IN ('failed','struggled','normal','strong')),
    journal                     TEXT,
    mistake_tags                JSONB NOT NULL DEFAULT '[]'::jsonb,
    leetcode_submission_id      TEXT
);

CREATE INDEX attempts_user_completed_idx
    ON attempts (user_id, completed_at DESC);

CREATE INDEX attempts_user_problem_idx
    ON attempts (user_id, problem_id, completed_at DESC);

CREATE UNIQUE INDEX attempts_user_lcid_uq
    ON attempts (user_id, leetcode_submission_id)
    WHERE leetcode_submission_id IS NOT NULL;

CREATE INDEX attempts_mistake_tags_gin ON attempts USING GIN (mistake_tags);

CREATE TABLE sessions (
    id                      BIGSERIAL PRIMARY KEY,
    user_id                 BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date                    DATE NOT NULL,
    problem_ids             JSONB NOT NULL DEFAULT '[]'::jsonb,
    completed_problem_ids   JSONB NOT NULL DEFAULT '[]'::jsonb,
    started_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at            TIMESTAMPTZ,
    UNIQUE (user_id, date)
);

CREATE TABLE auth_sessions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind        TEXT NOT NULL DEFAULT 'web' CHECK (kind IN ('web','ext')),
    token_hash  BYTEA UNIQUE NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX auth_sessions_user_idx ON auth_sessions (user_id, kind);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth_sessions;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS attempts;
DROP TABLE IF EXISTS user_problems;
DROP TABLE IF EXISTS problem_patterns;
DROP TABLE IF EXISTS patterns;
DROP TABLE IF EXISTS problems;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
