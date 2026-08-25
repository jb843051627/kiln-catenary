package store

const schemaSQL = `
CREATE TABLE IF NOT EXISTS kilns (
 id TEXT PRIMARY KEY, code TEXT NOT NULL UNIQUE, cell TEXT NOT NULL,
 max_temperature REAL NOT NULL, max_pressure REAL NOT NULL, atmosphere TEXT NOT NULL,
 active INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS zones (
 id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL REFERENCES kilns(id), name TEXT NOT NULL,
 target REAL NOT NULL, ramp_rate REAL NOT NULL, deviation_limit REAL NOT NULL, enabled INTEGER NOT NULL DEFAULT 1
);
CREATE TABLE IF NOT EXISTS thermal_stages (
 id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL REFERENCES kilns(id), name TEXT NOT NULL,
 kind TEXT NOT NULL, sequence INTEGER NOT NULL, start_temp REAL NOT NULL, end_temp REAL NOT NULL,
 hold_ns INTEGER NOT NULL, status TEXT NOT NULL, interlock TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS stages_kiln_sequence ON thermal_stages(kiln_id, sequence);
CREATE TABLE IF NOT EXISTS firing_runs (
 id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL REFERENCES kilns(id), status TEXT NOT NULL,
 started_at TEXT NOT NULL, finished_at TEXT NOT NULL, summary TEXT NOT NULL, score REAL NOT NULL
);
CREATE TABLE IF NOT EXISTS atmosphere_samples (
 id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL REFERENCES kilns(id), run_id TEXT NOT NULL REFERENCES firing_runs(id),
 temperature REAL NOT NULL, pressure REAL NOT NULL, oxidation REAL NOT NULL,
 observed_at TEXT NOT NULL, sequence INTEGER NOT NULL, source TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS samples_run_time ON atmosphere_samples(run_id, observed_at);
CREATE TABLE IF NOT EXISTS safety_events (
 id TEXT PRIMARY KEY, run_id TEXT NOT NULL, kiln_id TEXT NOT NULL REFERENCES kilns(id),
 kind TEXT NOT NULL, severity TEXT NOT NULL, message TEXT NOT NULL, created_at TEXT NOT NULL, resolved INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS events_kiln_time ON safety_events(kiln_id, created_at);
`
