ALTER TABLE nodes DROP COLUMN protocol_config;
ALTER TABLE servers ADD COLUMN protocols TEXT NULL;