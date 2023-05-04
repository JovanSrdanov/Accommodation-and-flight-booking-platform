-- public.account_credentials definition

-- Drop table

-- DROP TABLE public.account_credentials;

CREATE TABLE public.account_credentials (
	id text NOT NULL,
	email text NULL,
	"password" text NULL,
	salt text NULL,
	"role" int2 NULL,
	CONSTRAINT account_credentials_email_key UNIQUE (email),
	CONSTRAINT account_credentials_pkey PRIMARY KEY (id)
);

INSERT INTO public.account_credentials (id,email,"password",salt,"role") VALUES
	 ('8591f03d-e7fa-11ed-853e-040e3c52dc2b','host@example.com','password123','salt',0),
	 ('8c41edfd-e7fa-11ed-853e-040e3c52dc2b','guest@example.com','password123','salt',1);
