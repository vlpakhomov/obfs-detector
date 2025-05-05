BEGIN;


CREATE TABLE IF NOT EXISTS blocked_ip_address
(
    PRIMARY KEY (address),
    address               TEXT,
    verdict TEXT,
    created_at         timestamptz        NOT NULL DEFAULT now()
);

COMMIT;