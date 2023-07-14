CREATE TABLE IF NOT EXISTS queue (
    id INTEGER,
    campaign_id  INTEGER,
    name String,
    description String,
    priority INTEGER,
    removed Bool,
    EventTime DateTime('Europe/Moscow')
  ) ENGINE = NATS
    SETTINGS nats_url = 'localhost:4444',
             nats_subjects = 'item',
             nats_format = 'JSONEachRow',
             nats_max_block_size = 10,
             date_time_input_format = 'best_effort';