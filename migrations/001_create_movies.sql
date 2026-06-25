CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    title VARCHAR(200) NOT NULL,
    year INTEGER,
    about TEXT,
    facts TEXT,
    category VARCHAR(20) CHECK (category IN ('русский', 'зарубежный')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_movies_category ON movies(category);
CREATE INDEX IF NOT EXISTS idx_movies_key ON movies(key);