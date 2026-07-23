ALTER TABLE nodes ADD COLUMN protocol_config TEXT NULL;
COMMENT ON COLUMN nodes.protocol_config IS 'Protocol Configuration JSON';
ALTER TABLE servers DROP COLUMN protocols;