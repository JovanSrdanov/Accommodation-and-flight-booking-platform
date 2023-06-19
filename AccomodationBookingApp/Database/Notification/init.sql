-- public.notification_consent definition

-- Drop table

-- DROP TABLE public.notification_consent;

CREATE TABLE public.notification_consent (
	user_profile_id text NOT NULL,
	request_made bool NULL,
	reservation_canceled bool NULL,
	host_rating_given bool NULL,
	accommodation_rating_given bool NULL,
	prominent_host bool NULL,
	host_responded bool NULL,
	CONSTRAINT notification_consent_pkey PRIMARY KEY (user_profile_id)
);

INSERT INTO public.notification_consent (user_profile_id,request_made,reservation_canceled,host_rating_given,accommodation_rating_given,prominent_host,host_responded) VALUES
   ('5fa599f7-ef60-11ed-895f-0242ac1c0006',true,true,true,true,true,true),
   ('65ac16f5-ef60-11ed-895f-0242ac1c0006',true,true,true,true,true,true),
   ('686bab3c-afca-428e-afaf-b2b42cd9c7df',true,true,true,true,true,true);
