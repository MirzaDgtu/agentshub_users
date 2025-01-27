CREATE TABLE IF NOT EXISTS Users (
    UserID SERIAL PRIMARY KEY,
    Email VARCHAR(100) NOT NULL UNIQUE,
    PasswordHash VARCHAR(255) NOT NULL,
    FirstName VARCHAR(50),
    MiddleName VARCHAR(50),
    LastName VARCHAR(50) Null,
    ProfileImageURL VARCHAR(255),
    SignIn BOOLEAN DEFAULT FALSE, -- Поле для отслеживания статуса входа
    IsBlocked BOOLEAN DEFAULT FALSE, -- Поле для блокировки пользователя
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_email_format CHECK (Email ~* '^[^@\\s]+@[^@\\s]+\\.[^@\\s]+$')
);

CREATE INDEX IF NOT EXISTS idx_email ON Users (Email);

CREATE TABLE IF NOT EXISTS Roles (
    ID SERIAL PRIMARY KEY,
    RoleName VARCHAR(50) NOT NULL UNIQUE,
    Description VARCHAR(255),
    Priority INT NULL -- Поле для указания приоритета роли
);

CREATE TABLE IF NOT EXISTS UserRoles (
    ID SERIAL PRIMARY KEY,
    UserID INT NOT NULL,
    RoleID INT NOT NULL,
    FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE,
    FOREIGN KEY (RoleID) REFERENCES Roles(ID) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_role_name ON Roles (RoleName);
CREATE INDEX IF NOT EXISTS idx_user_role ON UserRoles (UserID, RoleID);

CREATE TABLE IF NOT EXISTS Permissions (
    ID SERIAL PRIMARY KEY,
    PermissionName VARCHAR(50) NOT NULL UNIQUE,
    Description VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS idx_permission_name ON Permissions (PermissionName);

