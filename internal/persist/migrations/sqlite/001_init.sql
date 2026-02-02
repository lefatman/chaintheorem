-- File: internal/persist/migrations/sqlite/001_init.sql
CREATE TABLE accounts (
	user_id INTEGER PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	username TEXT NOT NULL UNIQUE,
	pass_hash BLOB NOT NULL,
	created_at INTEGER NOT NULL,
	last_login_at INTEGER
);

CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	expires_at INTEGER NOT NULL,
	created_at INTEGER NOT NULL
);

CREATE TABLE army_loadouts (
	user_id INTEGER PRIMARY KEY,
	element_id INTEGER NOT NULL,
	army_ability_1 INTEGER NOT NULL DEFAULT 0,
	army_ability_2 INTEGER NOT NULL DEFAULT 0,
	army_ability_3 INTEGER NOT NULL DEFAULT 0,
	army_ability_4 INTEGER NOT NULL DEFAULT 0,
	ability_pawn INTEGER NOT NULL DEFAULT 0,
	ability_knight INTEGER NOT NULL DEFAULT 0,
	ability_bishop INTEGER NOT NULL DEFAULT 0,
	ability_rook INTEGER NOT NULL DEFAULT 0,
	ability_queen INTEGER NOT NULL DEFAULT 0,
	ability_king INTEGER NOT NULL DEFAULT 0,
	item_1 INTEGER NOT NULL DEFAULT 0,
	item_2 INTEGER NOT NULL DEFAULT 0,
	item_3 INTEGER NOT NULL DEFAULT 0,
	item_4 INTEGER NOT NULL DEFAULT 0,
	updated_at INTEGER NOT NULL
);

CREATE TABLE progression (
	user_id INTEGER PRIMARY KEY,
	level INTEGER NOT NULL,
	xp INTEGER NOT NULL
);

CREATE TABLE user_unlocks (
	user_id INTEGER NOT NULL,
	flag_id INTEGER NOT NULL,
	unlocked_at INTEGER NOT NULL,
	PRIMARY KEY (user_id, flag_id)
);
