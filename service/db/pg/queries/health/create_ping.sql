INSERT INTO pings (message)
VALUES ($1)
RETURNING message, EXTRACT(EPOCH FROM received_at)::bigint AS received_at_unix;
