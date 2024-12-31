CREATE TABLE IF NOT EXISTS operations (
	variable TEXT PRIMARY KEY,
	content TEXT
);

--this is for backing up past models
CREATE TABLE IF NOT EXISTS releases (
    id TEXT PRIMARY KEY,
    title TEXT,
    created_at DATE,
    content TEXT
);

CREATE TABLE IF NOT EXISTS controllers (
    id TEXT PRIMARY KEY,
    ip TEXT,
    mac TEXT,
    last_register DATE,
    rbmq_user TEXT,
    register_key TEXT,
    banned BOOLEAN,
    nodes_connected INT
);