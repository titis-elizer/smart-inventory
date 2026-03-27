-- DROP SCHEMA test_326;

CREATE SCHEMA test_326 AUTHORIZATION postgres;
-- test_326.audit_logs definition

-- Drop table

-- DROP TABLE test_326.audit_logs;

CREATE TABLE test_326.audit_logs (
	id uuid DEFAULT test_326.gen_random_uuid() NOT NULL,
	entity text NULL,
	entity_id uuid NULL,
	"action" text NULL,
	"data" jsonb NULL,
	created_at timestamp DEFAULT now() NULL,
	CONSTRAINT audit_logs_pkey PRIMARY KEY (id)
);


-- test_326.inventory_items definition

-- Drop table

-- DROP TABLE test_326.inventory_items;

CREATE TABLE test_326.inventory_items (
	id uuid DEFAULT test_326.gen_random_uuid() NOT NULL,
	"name" text NOT NULL,
	sku text NOT NULL,
	customer text NULL,
	physical_stock int4 DEFAULT 0 NOT NULL,
	reserved_stock int4 DEFAULT 0 NOT NULL,
	created_at timestamp DEFAULT now() NULL,
	updated_at timestamp DEFAULT now() NULL,
	CONSTRAINT check_available_stock CHECK (((physical_stock - reserved_stock) >= 0)),
	CONSTRAINT inventory_items_physical_stock_check CHECK ((physical_stock >= 0)),
	CONSTRAINT inventory_items_pkey PRIMARY KEY (id),
	CONSTRAINT inventory_items_reserved_stock_check CHECK ((reserved_stock >= 0)),
	CONSTRAINT inventory_items_sku_key UNIQUE (sku)
);
CREATE INDEX idx_inventory_customer ON test_326.inventory_items USING btree (customer);

-- Table Triggers

create trigger set_timestamp_inventory before
update
    on
    test_326.inventory_items for each row execute function test_326.update_updated_at_column();


-- test_326.stock_in definition

-- Drop table

-- DROP TABLE test_326.stock_in;

CREATE TABLE test_326.stock_in (
	id uuid DEFAULT test_326.gen_random_uuid() NOT NULL,
	status text NOT NULL,
	created_at timestamp DEFAULT now() NULL,
	updated_at timestamp DEFAULT now() NULL,
	done_at timestamp NULL,
	CONSTRAINT stock_in_pkey PRIMARY KEY (id),
	CONSTRAINT stock_in_status_check CHECK ((status = ANY (ARRAY['created'::text, 'in_progress'::text, 'done'::text, 'canceled'::text])))
);

-- Table Triggers

create trigger set_timestamp_stock_in before
update
    on
    test_326.stock_in for each row execute function test_326.update_updated_at_column();


-- test_326.stock_out definition

-- Drop table

-- DROP TABLE test_326.stock_out;

CREATE TABLE test_326.stock_out (
	id uuid DEFAULT test_326.gen_random_uuid() NOT NULL,
	status text NOT NULL,
	created_at timestamp DEFAULT now() NULL,
	updated_at timestamp DEFAULT now() NULL,
	done_at timestamp NULL,
	CONSTRAINT stock_out_pkey PRIMARY KEY (id),
	CONSTRAINT stock_out_status_check CHECK ((status = ANY (ARRAY['draft'::text, 'allocated'::text, 'in_progress'::text, 'done'::text, 'canceled'::text])))
);

-- Table Triggers

create trigger set_timestamp_stock_out before
update
    on
    test_326.stock_out for each row execute function test_326.update_updated_at_column();


-- test_326.stock_adjustments definition

-- Drop table

-- DROP TABLE test_326.stock_adjustments;

CREATE TABLE test_326.stock_adjustments (
	id uuid DEFAULT test_326.gen_random_uuid() NOT NULL,
	inventory_item_id uuid NULL,
	adjustment int4 NOT NULL,
	reason text NULL,
	created_at timestamp DEFAULT now() NULL,
	CONSTRAINT stock_adjustments_pkey PRIMARY KEY (id),
	CONSTRAINT stock_adjustments_inventory_item_id_fkey FOREIGN KEY (inventory_item_id) REFERENCES test_326.inventory_items(id)
);


-- test_326.stock_in_items definition

-- Drop table

-- DROP TABLE test_326.stock_in_items;

CREATE TABLE test_326.stock_in_items (
	id uuid DEFAULT test_326.gen_random_uuid() NOT NULL,
	stock_in_id uuid NOT NULL,
	inventory_item_id uuid NOT NULL,
	qty int4 NOT NULL,
	CONSTRAINT stock_in_items_pkey PRIMARY KEY (id),
	CONSTRAINT stock_in_items_qty_check CHECK ((qty > 0)),
	CONSTRAINT stock_in_items_inventory_item_id_fkey FOREIGN KEY (inventory_item_id) REFERENCES test_326.inventory_items(id),
	CONSTRAINT stock_in_items_stock_in_id_fkey FOREIGN KEY (stock_in_id) REFERENCES test_326.stock_in(id) ON DELETE CASCADE
);


-- test_326.stock_in_logs definition

-- Drop table

-- DROP TABLE test_326.stock_in_logs;

CREATE TABLE test_326.stock_in_logs (
	id uuid DEFAULT test_326.gen_random_uuid() NOT NULL,
	stock_in_id uuid NULL,
	status text NULL,
	note text NULL,
	created_at timestamp DEFAULT now() NULL,
	CONSTRAINT stock_in_logs_pkey PRIMARY KEY (id),
	CONSTRAINT stock_in_logs_stock_in_id_fkey FOREIGN KEY (stock_in_id) REFERENCES test_326.stock_in(id) ON DELETE CASCADE
);


-- test_326.stock_out_items definition

-- Drop table

-- DROP TABLE test_326.stock_out_items;

CREATE TABLE test_326.stock_out_items (
	id uuid DEFAULT test_326.gen_random_uuid() NOT NULL,
	stock_out_id uuid NULL,
	inventory_item_id uuid NULL,
	qty int4 NOT NULL,
	CONSTRAINT stock_out_items_pkey PRIMARY KEY (id),
	CONSTRAINT stock_out_items_qty_check CHECK ((qty > 0)),
	CONSTRAINT stock_out_items_inventory_item_id_fkey FOREIGN KEY (inventory_item_id) REFERENCES test_326.inventory_items(id),
	CONSTRAINT stock_out_items_stock_out_id_fkey FOREIGN KEY (stock_out_id) REFERENCES test_326.stock_out(id) ON DELETE CASCADE
);



-- DROP FUNCTION test_326.armor(bytea);

CREATE OR REPLACE FUNCTION test_326.armor(bytea)
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_armor$function$
;

-- DROP FUNCTION test_326.armor(bytea, _text, _text);

CREATE OR REPLACE FUNCTION test_326.armor(bytea, text[], text[])
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_armor$function$
;

-- DROP FUNCTION test_326.crypt(text, text);

CREATE OR REPLACE FUNCTION test_326.crypt(text, text)
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_crypt$function$
;

-- DROP FUNCTION test_326.dearmor(text);

CREATE OR REPLACE FUNCTION test_326.dearmor(text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_dearmor$function$
;

-- DROP FUNCTION test_326.decrypt(bytea, bytea, text);

CREATE OR REPLACE FUNCTION test_326.decrypt(bytea, bytea, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_decrypt$function$
;

-- DROP FUNCTION test_326.decrypt_iv(bytea, bytea, bytea, text);

CREATE OR REPLACE FUNCTION test_326.decrypt_iv(bytea, bytea, bytea, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_decrypt_iv$function$
;

-- DROP FUNCTION test_326.digest(text, text);

CREATE OR REPLACE FUNCTION test_326.digest(text, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_digest$function$
;

-- DROP FUNCTION test_326.digest(bytea, text);

CREATE OR REPLACE FUNCTION test_326.digest(bytea, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_digest$function$
;

-- DROP FUNCTION test_326.encrypt(bytea, bytea, text);

CREATE OR REPLACE FUNCTION test_326.encrypt(bytea, bytea, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_encrypt$function$
;

-- DROP FUNCTION test_326.encrypt_iv(bytea, bytea, bytea, text);

CREATE OR REPLACE FUNCTION test_326.encrypt_iv(bytea, bytea, bytea, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_encrypt_iv$function$
;

-- DROP FUNCTION test_326.gen_random_bytes(int4);

CREATE OR REPLACE FUNCTION test_326.gen_random_bytes(integer)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_random_bytes$function$
;

-- DROP FUNCTION test_326.gen_random_uuid();

CREATE OR REPLACE FUNCTION test_326.gen_random_uuid()
 RETURNS uuid
 LANGUAGE c
 PARALLEL SAFE
AS '$libdir/pgcrypto', $function$pg_random_uuid$function$
;

-- DROP FUNCTION test_326.gen_salt(text, int4);

CREATE OR REPLACE FUNCTION test_326.gen_salt(text, integer)
 RETURNS text
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_gen_salt_rounds$function$
;

-- DROP FUNCTION test_326.gen_salt(text);

CREATE OR REPLACE FUNCTION test_326.gen_salt(text)
 RETURNS text
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_gen_salt$function$
;

-- DROP FUNCTION test_326.hmac(text, text, text);

CREATE OR REPLACE FUNCTION test_326.hmac(text, text, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_hmac$function$
;

-- DROP FUNCTION test_326.hmac(bytea, bytea, text);

CREATE OR REPLACE FUNCTION test_326.hmac(bytea, bytea, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pg_hmac$function$
;

-- DROP FUNCTION test_326.pgp_armor_headers(in text, out text, out text);

CREATE OR REPLACE FUNCTION test_326.pgp_armor_headers(text, OUT key text, OUT value text)
 RETURNS SETOF record
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_armor_headers$function$
;

-- DROP FUNCTION test_326.pgp_key_id(bytea);

CREATE OR REPLACE FUNCTION test_326.pgp_key_id(bytea)
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_key_id_w$function$
;

-- DROP FUNCTION test_326.pgp_pub_decrypt(bytea, bytea, text);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_decrypt(bytea, bytea, text)
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_decrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_pub_decrypt(bytea, bytea);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_decrypt(bytea, bytea)
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_decrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_pub_decrypt(bytea, bytea, text, text);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_decrypt(bytea, bytea, text, text)
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_decrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_pub_decrypt_bytea(bytea, bytea);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_decrypt_bytea(bytea, bytea)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_decrypt_bytea$function$
;

-- DROP FUNCTION test_326.pgp_pub_decrypt_bytea(bytea, bytea, text);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_decrypt_bytea(bytea, bytea, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_decrypt_bytea$function$
;

-- DROP FUNCTION test_326.pgp_pub_decrypt_bytea(bytea, bytea, text, text);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_decrypt_bytea(bytea, bytea, text, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_decrypt_bytea$function$
;

-- DROP FUNCTION test_326.pgp_pub_encrypt(text, bytea);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_encrypt(text, bytea)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_encrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_pub_encrypt(text, bytea, text);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_encrypt(text, bytea, text)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_encrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_pub_encrypt_bytea(bytea, bytea, text);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_encrypt_bytea(bytea, bytea, text)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_encrypt_bytea$function$
;

-- DROP FUNCTION test_326.pgp_pub_encrypt_bytea(bytea, bytea);

CREATE OR REPLACE FUNCTION test_326.pgp_pub_encrypt_bytea(bytea, bytea)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_pub_encrypt_bytea$function$
;

-- DROP FUNCTION test_326.pgp_sym_decrypt(bytea, text, text);

CREATE OR REPLACE FUNCTION test_326.pgp_sym_decrypt(bytea, text, text)
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_sym_decrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_sym_decrypt(bytea, text);

CREATE OR REPLACE FUNCTION test_326.pgp_sym_decrypt(bytea, text)
 RETURNS text
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_sym_decrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_sym_decrypt_bytea(bytea, text, text);

CREATE OR REPLACE FUNCTION test_326.pgp_sym_decrypt_bytea(bytea, text, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_sym_decrypt_bytea$function$
;

-- DROP FUNCTION test_326.pgp_sym_decrypt_bytea(bytea, text);

CREATE OR REPLACE FUNCTION test_326.pgp_sym_decrypt_bytea(bytea, text)
 RETURNS bytea
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_sym_decrypt_bytea$function$
;

-- DROP FUNCTION test_326.pgp_sym_encrypt(text, text, text);

CREATE OR REPLACE FUNCTION test_326.pgp_sym_encrypt(text, text, text)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_sym_encrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_sym_encrypt(text, text);

CREATE OR REPLACE FUNCTION test_326.pgp_sym_encrypt(text, text)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_sym_encrypt_text$function$
;

-- DROP FUNCTION test_326.pgp_sym_encrypt_bytea(bytea, text, text);

CREATE OR REPLACE FUNCTION test_326.pgp_sym_encrypt_bytea(bytea, text, text)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_sym_encrypt_bytea$function$
;

-- DROP FUNCTION test_326.pgp_sym_encrypt_bytea(bytea, text);

CREATE OR REPLACE FUNCTION test_326.pgp_sym_encrypt_bytea(bytea, text)
 RETURNS bytea
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/pgcrypto', $function$pgp_sym_encrypt_bytea$function$
;

-- DROP FUNCTION test_326.update_timestamp();

CREATE OR REPLACE FUNCTION test_326.update_timestamp()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$function$
;

-- DROP FUNCTION test_326.update_updated_at_column();

CREATE OR REPLACE FUNCTION test_326.update_updated_at_column()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$function$
;

-- DROP FUNCTION test_326.uuid_generate_v1();

CREATE OR REPLACE FUNCTION test_326.uuid_generate_v1()
 RETURNS uuid
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v1$function$
;

-- DROP FUNCTION test_326.uuid_generate_v1mc();

CREATE OR REPLACE FUNCTION test_326.uuid_generate_v1mc()
 RETURNS uuid
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v1mc$function$
;

-- DROP FUNCTION test_326.uuid_generate_v3(uuid, text);

CREATE OR REPLACE FUNCTION test_326.uuid_generate_v3(namespace uuid, name text)
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v3$function$
;

-- DROP FUNCTION test_326.uuid_generate_v4();

CREATE OR REPLACE FUNCTION test_326.uuid_generate_v4()
 RETURNS uuid
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v4$function$
;

-- DROP FUNCTION test_326.uuid_generate_v5(uuid, text);

CREATE OR REPLACE FUNCTION test_326.uuid_generate_v5(namespace uuid, name text)
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v5$function$
;

-- DROP FUNCTION test_326.uuid_nil();

CREATE OR REPLACE FUNCTION test_326.uuid_nil()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_nil$function$
;

-- DROP FUNCTION test_326.uuid_ns_dns();

CREATE OR REPLACE FUNCTION test_326.uuid_ns_dns()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_ns_dns$function$
;

-- DROP FUNCTION test_326.uuid_ns_oid();

CREATE OR REPLACE FUNCTION test_326.uuid_ns_oid()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_ns_oid$function$
;

-- DROP FUNCTION test_326.uuid_ns_url();

CREATE OR REPLACE FUNCTION test_326.uuid_ns_url()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_ns_url$function$
;

-- DROP FUNCTION test_326.uuid_ns_x500();

CREATE OR REPLACE FUNCTION test_326.uuid_ns_x500()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_ns_x500$function$
;