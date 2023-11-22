BEGIN;

CREATE TABLE users (
    id serial,
    name varchar(80) NOT NULL,
    login varchar(100) NOT NULL UNIQUE,
    password varchar(200) NOT NULL,
    created_at timestamp DEFAULT current_timestamp,
    modified_at timestamp NOT NULL,
    deleted bool NOT NULL DEFAULT false,
    last_login timestamp DEFAULT current_timestamp,

    PRIMARY KEY(id)
);

COMMIT;