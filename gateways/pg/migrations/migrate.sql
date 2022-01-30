-- This is a temporary file to create the schemas, not a real migration.

CREATE DATABASE library;

CREATE TABLE shelves (
    name TEXT PRIMARY KEY,
    create_time TIMESTAMP,
    update_time TIMESTAMP
);

CREATE TABLE books (
    name TEXT,
    author TEXT,
    shelf_name TEXT,
    create_time TIMESTAMP,
    update_time TIMESTAMP,
    CONSTRAINT fk_shelf FOREIGN KEY (shelf_name) REFERENCES shelves (name),
    PRIMARY KEY (name, shelf_name)
);