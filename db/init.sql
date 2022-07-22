CREATE TABLE shortURL (
    hashcode varchar(10) NOT NULL UNIQUE PRIMARY KEY,
    link varchar(250) NOT NULL
);