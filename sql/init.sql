CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE notifications (
    id uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    chatId INTEGER NOT NULL ,
    message VARCHAR(100) NOT NULL,
    notificationDate DATE NOT NULL
)