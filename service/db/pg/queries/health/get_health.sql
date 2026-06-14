SELECT 'ok'::text AS status, EXTRACT(EPOCH FROM NOW())::bigint AS checked_at_unix;
