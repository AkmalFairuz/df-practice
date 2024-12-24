CREATE TABLE users
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    displayName  VARCHAR(255) NOT NULL,
    rankName     VARCHAR(255) NOT NULL DEFAULT 'default',
    rankExpireAt BIGINT       NOT NULL DEFAULT 0,
    xuid         VARCHAR(255) NOT NULL,
    exp          BIGINT       NOT NULL DEFAULT 0,
    registeredAt BIGINT       NOT NULL,
    lastSeenAt   BIGINT       NOT NULL,
    UNIQUE (xuid)
) ENGINE = InnoDB;