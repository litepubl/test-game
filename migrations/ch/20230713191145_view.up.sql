  CREATE MATERIALIZED VIEW consumer TO items
    AS SELECT id, campaign_id, name, description, priority, removed, EventTime FROM queue;