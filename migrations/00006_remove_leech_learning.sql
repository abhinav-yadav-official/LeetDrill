-- +goose Up
-- +goose StatementBegin

UPDATE user_problems SET status = 'review' WHERE status IN ('leech', 'learning');

ALTER TABLE user_problems DROP CONSTRAINT IF EXISTS user_problems_status_check;
ALTER TABLE user_problems ADD CONSTRAINT user_problems_status_check
    CHECK (status IN ('new', 'review', 'mastered'));

DROP INDEX IF EXISTS user_problems_due_idx;
CREATE INDEX user_problems_due_idx ON user_problems (user_id, next_due_at)
    WHERE status = 'review';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_problems DROP CONSTRAINT IF EXISTS user_problems_status_check;
ALTER TABLE user_problems ADD CONSTRAINT user_problems_status_check
    CHECK (status IN ('new','learning','review','mastered','leech'));

DROP INDEX IF EXISTS user_problems_due_idx;
CREATE INDEX user_problems_due_idx ON user_problems (user_id, next_due_at)
    WHERE status NOT IN ('leech','new');
-- +goose StatementEnd
