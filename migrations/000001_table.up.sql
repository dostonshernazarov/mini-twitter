CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    username VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL,
    bio TEXT,
    profile_picture TEXT,
    password TEXT NOT NULL,
    role VARCHAR(10) NOT NULL,
    refresh TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tweets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    parent_tweet_id UUID,
    content VARCHAR(280),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (parent_tweet_id) REFERENCES tweets(id)
);

CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY,
    tweet_id UUID NOT NULL,
    file_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (tweet_id) REFERENCES tweets(id)
);

CREATE TABLE IF NOT EXISTS follows (
    user_id UUID NOT NULL,
    following_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (following_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS likes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    tweet_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (tweet_id) REFERENCES tweets(id)
);


INSERT INTO users (id, name, username, email, password, role) VALUES
('1a8d3b24-3c29-4d21-a7d6-fcdd5a92c56a', 'Jhon Doe', 'jhon_doe', 'adminemail@gmail.com', '$2a$14$0JDtvVEXdegqfz/Q.ThLm.Hg4kes50BRkBPBI48DbDKiI0Z9ifE9O', 'admin');