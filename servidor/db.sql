USE sociable;

CREATE TABLE IF NOT EXISTS StatusList (
    ID INT NOT NULL AUTO_INCREMENT,
    StatusMessage TEXT(15),
    PRIMARY KEY (ID)
);

CREATE TABLE IF NOT EXISTS Users (
    ID INT NOT NULL AUTO_INCREMENT,
    Username TEXT(20),
    Pass TEXT(64),
    Email TEXT(320),
    EmailCode TEXT(64),
    Avatar TEXT(64),
    Admin BOOLEAN,
    Banned BOOLEAN,
    Verified BOOLEAN,
    StatusID INT,
    PRIMARY KEY (ID),
    FOREIGN KEY (StatusID) REFERENCES StatusList(ID)
);

CREATE TABLE IF NOT EXISTS Blocks (
    BlockTo INT,
    BlockBy INT,
    FOREIGN KEY (BlockTo) REFERENCES Users(ID),
    FOREIGN KEY (BlockBy) REFERENCES Users(ID)
);

CREATE TABLE IF NOT EXISTS Chats (
    UserID1 INT,
    UserID2 INT,
    User1 BOOLEAN,
    User2 BOOLEAN,
    ID INT NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (ID),
    FOREIGN KEY (UserID1) REFERENCES Users(ID),
    FOREIGN KEY (UserID2) REFERENCES Users(ID)
);

CREATE TABLE IF NOT EXISTS Messages (
    ID INT NOT NULL AUTO_INCREMENT,
    ChatID INT,
    AuthorID INT,
    Content TEXT(1000),
    Timestamp int,
    PRIMARY KEY (ID),
    FOREIGN KEY (ChatID) REFERENCES Chats(ID),
    FOREIGN KEY (AuthorID) REFERENCES Users(ID)
);

CREATE TABLE IF NOT EXISTS Posts (
    ID INT NOT NULL AUTO_INCREMENT,
    AuthorID INT,
    HiddenPost BOOLEAN,
    Title TEXT(30),
    Content TEXT(10000),
    PRIMARY KEY (ID),
    FOREIGN KEY (AuthorID) REFERENCES Users(ID)
);

INSERT INTO StatusList SET StatusMessage = "Desconectado";
INSERT INTO StatusList SET StatusMessage = "Conectado";