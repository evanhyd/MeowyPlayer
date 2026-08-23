PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS user_profile (
    user_id TEXT NOT NULL,
    username TEXT NOT NULL,
    language INTEGER NOT NULL,
    registration_date INTEGER NOT NULL, -- Unix
    token TEXT NOT NULL,
    PRIMARY KEY(user_id)
);

CREATE TABLE IF NOT EXISTS playlist (
    user_id TEXT NOT NULL,
    playlist_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    modified_date INTEGER NOT NULL,  -- Unix nano
    cover_blob BLOB NOT NULL,

    PRIMARY KEY(user_id, playlist_id)
);

CREATE TABLE IF NOT EXISTS music (
    music_id TEXT NOT NULL,
    source INTEGER NOT NULL,
    title TEXT NOT NULL,
    length_seconds INTEGER NOT NULL,
    
    PRIMARY KEY(music_id, source)
);

CREATE TABLE IF NOT EXISTS playlist_music (
    user_id TEXT,
    playlist_id INTEGER,
    music_id TEXT,
    source INTEGER,
    added_at INTEGER NOT NULL,      -- Unix nano

    PRIMARY KEY(user_id, playlist_id, music_id, source),
    FOREIGN KEY(user_id, playlist_id) REFERENCES playlist(user_id, playlist_id) ON DELETE CASCADE,
    FOREIGN KEY(music_id, source) REFERENCES music(music_id, source) ON DELETE RESTRICT
);
