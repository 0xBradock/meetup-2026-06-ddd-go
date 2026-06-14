IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'pings')
BEGIN
    CREATE TABLE pings (
        id          BIGINT IDENTITY(1,1) PRIMARY KEY,
        message     NVARCHAR(MAX)        NOT NULL,
        received_at DATETIMEOFFSET       NOT NULL DEFAULT GETUTCDATE()
    );
END
