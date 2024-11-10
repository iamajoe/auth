-- REF: iamajoe - for migrations to run we need to prepare things

CREATE SCHEMA IF NOT EXISTS {{ index .Options "Namespace" }};

DO $$
BEGIN
  -- TODO: do we like this on the migrations? is this even working?
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '{{ index .Options "Username" }}') THEN
      CREATE ROLE {{ index .Options "Username" }};
   END IF;
END $$;

do $$ begin
    create type {{ index .Options "Namespace" }}.factor_type as enum('totp', 'webauthn');
    create type {{ index .Options "Namespace" }}.factor_status as enum('unverified', 'verified');
    create type {{ index .Options "Namespace" }}.aal_level as enum('aal1', 'aal2', 'aal3');
exception
    when duplicate_object then null;
end $$;
