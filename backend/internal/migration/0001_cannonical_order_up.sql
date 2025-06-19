DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type')
        THEN CREATE TYPE status_type AS ENUM ('active', 'deprecated');
    END IF;
END $$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS cannonical_order (
    id BIGSERIAL PRIMARY KEY,
    real_vibe VARCHAR(255) NOT NULL UNIQUE,
    vibe_order BIGINT NOT NULL UNIQUE,
    version VARCHAR(20),
    status status_type NOT NULL DEFAULT 'active'
)