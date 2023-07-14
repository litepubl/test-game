CREATE TABLE IF NOT EXISTS items (
    id INTEGER,
    campaign_id  INTEGER,
    name String,
    description String,
    priority INTEGER,
    removed Bool,
    EventTime DateTime('Europe/Moscow')
)
    ENGINE = MergeTree() ORDER BY id;