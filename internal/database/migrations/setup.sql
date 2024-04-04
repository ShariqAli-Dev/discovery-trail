CREATE TABLE users (
    id  VARCHAR NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL
);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
