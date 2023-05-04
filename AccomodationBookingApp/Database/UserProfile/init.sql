-- public.address definition

-- Drop table

-- DROP TABLE public.address;

CREATE TABLE public.address (
	id text NOT NULL,
	country text NULL,
	city text NULL,
	street text NULL,
	street_number text NULL,
	CONSTRAINT address_pkey PRIMARY KEY (id)
);


-- public.user_profile definition

-- Drop table

-- DROP TABLE public.user_profile;

CREATE TABLE public.user_profile (
	id text NOT NULL,
	"name" text NULL,
	surname text NULL,
	email text NULL,
	address_id text NULL,
	CONSTRAINT user_profile_pkey PRIMARY KEY (id),
	CONSTRAINT fk_user_profile_address FOREIGN KEY (address_id) REFERENCES public.address(id)
);

INSERT INTO public.address (id,country,city,street,street_number) VALUES
	 ('0171c17d-ea79-11ed-b73d-040e3c52dc2b','USA','New York','Broadway','123'),
	 ('07394a7a-ea79-11ed-b73d-040e3c52dc2b','Canada','Toronto','Yonge Street','456');

INSERT INTO public.user_profile (id,"name",surname,email,address_id) VALUES
	 ('017182a4-ea79-11ed-b73d-040e3c52dc2b','John','Doe','johndoe@example.com','0171c17d-ea79-11ed-b73d-040e3c52dc2b'),
	 ('07394a37-ea79-11ed-b73d-040e3c52dc2b','Jane','Smith','janesmith@example.com','07394a7a-ea79-11ed-b73d-040e3c52dc2b');

