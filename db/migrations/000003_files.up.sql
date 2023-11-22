BEGIN;

CREATE TABLE files (
    id serial,
    folder_id int,
    owner_id int,
    name varchar(200) NOT NULL,
    type varchar(50) NOT NULL,
    path varchar(250) NOT NULL,
    created_at timestamp DEFAULT current_timestamp,
    modified_at timestamp NOT NULL,
    deleted bool NOT NULL DEFAULT false,

    PRIMARY KEY (id),
    CONSTRAINT fk_folders
        FOREIGN KEY(folder_id)
        REFERENCES folders(id),
    CONSTRAINT fk_owner
        FOREIGN KEY(owner_id)
        REFERENCES users(id)
);

COMMIT;