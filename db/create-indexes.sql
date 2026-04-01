CREATE INDEX IF NOT EXISTS object_type_hash_idx ON objects(type, hash) CLUSTER;
