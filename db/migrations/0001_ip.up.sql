CREATE TABLE ips (
  id serial primary key,
  account_id varchar(255),
  created_at timestamptz DEFAULT current_timestamp,
  updated_at timestamptz DEFAULT NULL,
  ip_address inet DEFAULT NULL
);

CREATE FUNCTION set_updated_at()
  RETURNS trigger as $$
  BEGIN
    NEW.updated_at := current_timestamp;
    RETURN NEW;
  END $$ language plpgsql;

CREATE TRIGGER ips_updated_at
  BEFORE UPDATE OR INSERT ON ips
  FOR EACH ROW
  EXECUTE PROCEDURE set_updated_at();