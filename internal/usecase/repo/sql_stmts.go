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
	// Deprecated:
	_upsertMetricOld = `
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
SELECT $1,$2,$3,$4,$5
WHERE NOT EXISTS (
    SELECT 1 FROM metrics WHERE id = $1
);
`
	_createNamedRow = `
INSERT INTO metrics (id, type, delta, value, hash)
SELECT :id,:type,:delta,:value,:hash
WHERE NOT EXISTS (
    SELECT 1 FROM metrics WHERE id = :id
);
`
	_upsertNamed = `
INSERT INTO metrics (id, type, delta, value, hash)
VALUES(:id, :type, :delta, :value, :hash) 
ON CONFLICT (id) DO UPDATE
SET delta = EXCLUDED.delta, value = EXCLUDED.value, hash = EXCLUDED.hash;
`
	_upsertGauge = `
INSERT INTO metrics (id, type, delta, value, hash)
VALUES(:id, :type, null, :value, :hash) 
ON CONFLICT (id) DO UPDATE
SET value = EXCLUDED.value, hash = EXCLUDED.hash;
`
	_upsertCounter = `
INSERT INTO metrics (id, type, delta, value, hash)
VALUES(:id, :type, :delta, null, :hash) 
ON CONFLICT (id) DO UPDATE
SET delta = metrics.delta + EXCLUDED.delta, hash = EXCLUDED.hash;
`
	_upsertMetric = `
INSERT INTO metrics (id, type, delta, value, hash)
VALUES (:id, :type, :delta, :value, :hash)
ON CONFLICT (id)
    DO UPDATE SET delta = CASE
                              WHEN metrics.type = 'counter'
                                  THEN metrics.delta + EXCLUDED.delta
                              END,
                  value = CASE
                              WHEN metrics.type = 'gauge'
                                  THEN EXCLUDED.value
                              END;
`
	_selectIds = `
SELECT id
FROM metrics;
`
)
