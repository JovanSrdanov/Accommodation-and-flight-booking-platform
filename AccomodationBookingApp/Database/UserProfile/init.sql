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
('5f9a433e-ef60-11ed-84fa-0242ac1c0007','Serbia','Belgrade','Kneza Milosa','50'),
('65a13dcf-ef60-11ed-84fa-0242ac1c0007','Serbia','Novi Sad','Bulevar Oslobođenja','115'),
('ae8f9903-c16c-4975-b250-aa924df11a61','Serbia','Novi Sad','Kralja petra','115');

INSERT INTO public.user_profile (id,"name",surname,email,address_id) VALUES
('5f9a3e86-ef60-11ed-84fa-0242ac1c0007','Marko','Maslesa','markomaslesa@gmail.com','5f9a433e-ef60-11ed-84fa-0242ac1c0007'),
('65a13dae-ef60-11ed-84fa-0242ac1c0007','Jana','Jankovic','janajankovic@gmail.com','65a13dcf-ef60-11ed-84fa-0242ac1c0007'),
('6198a15e-5751-4252-b15f-d5b01813dc15','Jovan','Srdanov','da@gmail.com','ae8f9903-c16c-4975-b250-aa924df11a61');

