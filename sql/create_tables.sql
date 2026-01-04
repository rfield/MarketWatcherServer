-- Table: public.accounts

-- DROP TABLE IF EXISTS public.accounts;

CREATE TABLE IF NOT EXISTS public.accounts
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name text COLLATE pg_catalog."default",
    user_id uuid NOT NULL,
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.accounts
    OWNER to postgres;

-- Table: public.assets

-- DROP TABLE IF EXISTS public.assets;

CREATE TABLE IF NOT EXISTS public.assets
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    account_id uuid NOT NULL,
    ticker text COLLATE pg_catalog."default" NOT NULL,
    holding_amount numeric,
    CONSTRAINT assets_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.assets
    OWNER to postgres;

-- Table: public.users

-- DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    username text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    token text COLLATE pg_catalog."default",
    given_name text COLLATE pg_catalog."default",
    family_name text COLLATE pg_catalog."default",
    email text COLLATE pg_catalog."default",
    CONSTRAINT user_pkey1 PRIMARY KEY (id),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_username_key UNIQUE (username)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;