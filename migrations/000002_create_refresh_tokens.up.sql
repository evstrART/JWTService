CREATE TABLE refresh_tokens (
                                id SERIAL PRIMARY KEY,
                                user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                token_id UUID NOT NULL UNIQUE,        -- jti из токена
                                expires_at TIMESTAMP NOT NULL,
                                revoked BOOLEAN DEFAULT FALSE,
                                created_at TIMESTAMP DEFAULT NOW()
);