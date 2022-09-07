package repo

const (
	_schema = `
CREATE TABLE IF NOT EXISTS metrics(
    "id" VARCHAR(255) UNIQUE NOT NULL,
    "type" VARCHAR(50) NOT NULL,
    "delta" BIGINT,
    "value" DOUBLE PRECISION,
    "hash" VARCHAR(255)
    );
`
	_upsertMetric = `
INSERT INTO metrics (id, type, delta, value, hash)
VALUES($1,$2,$3,$4,$5) 
ON CONFLICT (id) DO UPDATE
SET delta = EXCLUDED.delta, value = EXCLUDED.value, hash = EXCLUDED.hash;
`
	_readMetrics = `SELECT * FROM metrics;`
	_readMetric  = `
SELECT *
FROM metrics
WHERE id = $1 AND type = $2;
`
	_createRow = `
INSERT INTO metrics (id, type, delta, value, hash)
SELECT $1::VARCHAR,$2,$3,$4,$5
WHERE NOT EXISTS (
    SELECT 1 FROM metrics WHERE id = $1
);
`
)
