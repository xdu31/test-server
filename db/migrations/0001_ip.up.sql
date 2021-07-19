CREATE TABLE ips (
  id integer NOT NULL,
  account_id integer,
  ip_address  varchar(255) NOT NULL,
  created_at timestamptz DEFAULT current_timestamp,
  updated_at timestamptz DEFAULT NULL
);

CREATE SEQUENCE ips_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE ips_id_seq OWNED BY ips.id;

ALTER TABLE ONLY ips ALTER COLUMN id SET DEFAULT nextval('ips_id_seq'::regclass);

ALTER TABLE ONLY ips
    ADD CONSTRAINT pk_ips PRIMARY KEY (id);