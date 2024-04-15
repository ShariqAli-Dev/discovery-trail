CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE TABLE accounts (
    id  VARCHAR NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    credits INT NOT NULL
);

CREATE TABLE courses (
    id INT PRIMARY KEY,
    name VARCHAR NOT NULL,
    image VARCHAR NOT NULL,
    account_id VARCHAR NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE units (
    id INT PRIMARY KEY,
    name VARCHAR NOT NULL,
    course_id INT NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses(id)
);

CREATE TABLE chapters (
    id INT PRIMARY KEY,
    name VARCHAR NOT NULL,
    youtubeSearchQuery VARCHAR NOT NULL,
    videoID VARCHAR NOT NULL,
    summary VARCHAR NOT NULL,
    unit_id INT NOT NULL,
    FOREIGN KEY (unit_id) REFERENCES units(id)
);

CREATE TABLE questions (
    id INT PRIMARY KEY,
    question VARCHAR NOT NULL,
    answer VARCHAR NOT NULL,
    options VARCHAR NOT NULL,
    chapter_id INT NOT NULL,
    FOREIGN KEY (chapter_id) REFERENCES chapters(id)
);