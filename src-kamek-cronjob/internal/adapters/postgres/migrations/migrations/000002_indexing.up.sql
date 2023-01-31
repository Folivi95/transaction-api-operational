CREATE INDEX merchant_id ON transactions_api USING btree ((transaction ->'card_acceptor'->>'id'));
CREATE INDEX transaction_id ON transactions_api USING btree ((transaction ->>'id'));
