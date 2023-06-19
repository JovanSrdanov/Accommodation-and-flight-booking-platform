-- public.account_credentials definition

-- Drop table

-- DROP TABLE public.account_credentials;

CREATE TABLE public.account_credentials (
	id text NOT NULL,
	username text NULL,
	"password" text NULL,
	"role" int2 NULL,
	user_profile_id text NULL,
	CONSTRAINT account_credentials_pkey PRIMARY KEY (id),
	CONSTRAINT account_credentials_username_key UNIQUE (username)
);



INSERT INTO public.account_credentials (id,username,"password","role",user_profile_id) VALUES
   ('5fa599f7-ef60-11ed-895f-0242ac1c0006','host','$2a$10$fuzgf1lY3Lp0EWA/DT55yewRoGOPdkhJ077.8.MSNTUnbmjs0KOnS',0,'5f9a3e86-ef60-11ed-84fa-0242ac1c0007'),
   ('65ac16f5-ef60-11ed-895f-0242ac1c0006','guest1','$2a$10$ZOdGfHondlDV9unUkDl5fupfkucY4xow8vteWMPaOE8RkTgx40MKq',1,'65a13dae-ef60-11ed-84fa-0242ac1c0007'),
   ('686bab3c-afca-428e-afaf-b2b42cd9c7df','guest2','$2a$10$ZOdGfHondlDV9unUkDl5fupfkucY4xow8vteWMPaOE8RkTgx40MKq',1,'6198a15e-5751-4252-b15f-d5b01813dc15');

