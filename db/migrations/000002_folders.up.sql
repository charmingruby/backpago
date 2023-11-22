BEGIN;

CREATE TABLE folders (
    id serial,
    parent_id int,
    name varchar(60) NOT NULL,
    created_at timestamp DEFAULT current_timestamp,
    modified_at timestamp NOT NULL,
    deleted bool NOT NULL DEFAULT false,

    PRIMARY KEY (id),
    CONSTRAINT fk_parent
        FOREIGN KEY(parent_id)
        REFERENCES folders(id)
);

COMMIT;