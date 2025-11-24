PRAGMA foreign_keys = ON;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    salt TEXT NOT NULL
);

-- User sessions table
CREATE TABLE IF NOT EXISTS user_sessions (
    user_id INTEGER PRIMARY KEY,
    token TEXT UNIQUE NOT NULL,
    created_at INTEGER NOT NULL,  -- Unix nano
    FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Playlist table
CREATE TABLE IF NOT EXISTS playlists (
    user_id INTEGER NOT NULL,
    playlist_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    modified_date INTEGER NOT NULL,  -- Unix nano
    cover_blob BLOB NOT NULL,
    PRIMARY KEY(user_id, playlist_id),
    FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Music table
CREATE TABLE IF NOT EXISTS music (
    music_id TEXT NOT NULL,
    source INTEGER NOT NULL,
    title TEXT NOT NULL,
    length_seconds INTEGER NOT NULL,
    PRIMARY KEY(music_id, source)
);

-- Playlist_Music table
CREATE TABLE IF NOT EXISTS playlist_music (
    user_id INTEGER NOT NULL,
    playlist_id INTEGER NOT NULL,
    music_id TEXT NOT NULL,
    source INTEGER NOT NULL,
    added_at INTEGER NOT NULL,      -- Unix nano
    PRIMARY KEY(playlist_id, music_id, source),
    FOREIGN KEY(user_id, playlist_id) REFERENCES playlists(user_id, playlist_id) ON DELETE CASCADE,
    FOREIGN KEY(music_id, source) REFERENCES music(music_id, source) ON DELETE CASCADE
);

INSERT OR IGNORE INTO users VALUES (0, 'default', '', '');
