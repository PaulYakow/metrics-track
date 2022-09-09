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
	_selectAllMetrics = `
SELECT *
FROM metrics;
`
	_selectMetricByIdAndType = `
SELECT *
FROM metrics
WHERE id = $1 AND type = $2;
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
)
