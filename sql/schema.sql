CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    type INTEGER NOT NULL,
    value INTEGER NOT NULL,
    state TEXT NOT NULL,
    creation_time TIMESTAMP NOT NULL,
    last_update_time TIMESTAMP NOT NULL
);
