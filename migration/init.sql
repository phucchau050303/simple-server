CREATE TABLE users (
    ID INT PRIMARY KEY,
    Email VARCHAR(255) NOT NULL UNIQUE,
    PasswordHash VARCHAR(255) NOT NULL,
    CreatedAt DATE,
    DeletedAt DATE
)