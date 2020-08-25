CREATE TABLE IF NOT EXISTS project
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	slug TEXT NOT NULL,
	description TEXT NOT NULL,
	tags TEXT,
	image TEXT,
	repo TEXT,
	demo TEXT,
	is_hidden INTEGER,
	added_on INTEGER,
	edited_on INTEGER
);