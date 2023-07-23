ALTER TABLE rainbow_chain ALTER COLUMN end_hash TYPE BYTEA USING decode(end_hash, 'hex');

ALTER TABLE rainbow_table_search ALTER COLUMN hash TYPE BYTEA USING decode(hash, 'hex');
