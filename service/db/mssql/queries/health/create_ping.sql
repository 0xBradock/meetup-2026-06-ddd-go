INSERT INTO pings (message)
OUTPUT INSERTED.message, DATEDIFF_BIG(SECOND, '19700101', INSERTED.received_at) AS received_at_unix
VALUES (@p1);
