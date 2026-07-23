ALTER TABLE nodes ADD COLUMN protocol_config TEXT NULL COMMENT 'Protocol Configuration JSON';
ALTER TABLE servers DROP COLUMN protocols;