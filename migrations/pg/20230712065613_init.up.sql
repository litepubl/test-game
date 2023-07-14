CREATE TABLE IF NOT EXISTS campaigns (
    id serial PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    id serial PRIMARY KEY,
    campaign_id  INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    priority INTEGER NOT NULL,
    removed bool NOT NULL default false,
    created_at timestamp(0) with time zone NOT NULL
);

CREATE INDEX idx_items_campaign_id  ON items (campaign_id );
CREATE INDEX idx_items_priority ON items (priority);

ALTER TABLE items ADD CONSTRAINT items_campaign_id_fk FOREIGN KEY (campaign_id ) REFERENCES campaigns(id) ON DELETE RESTRICT NOT DEFERRABLE INITIALLY IMMEDIATE;

INSERT INTO campaigns  (name) values ('Первая запись');
