-- Table: public.satellites

-- DROP TABLE IF EXISTS public.satellites;

CREATE TABLE IF NOT EXISTS public.satellites
(
    id uuid NOT NULL,
    name character varying(20) COLLATE pg_catalog."default" NOT NULL,
    "order" "char",
    position_id uuid,
    distance double precision,
    message text[] COLLATE pg_catalog."default",
    CONSTRAINT "Satellites_pkey" PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.satellites
    OWNER to postgres;

-- Table: public.positions

-- DROP TABLE IF EXISTS public.positions;

CREATE TABLE IF NOT EXISTS public.positions
(
    id uuid NOT NULL,
    x double precision,
    y double precision,
    CONSTRAINT "Positions_pkey" PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.positions
    OWNER to postgres;

INSERT INTO public.satellites(id, name, "order", position_id)
	VALUES ('1c5867bb-6c45-4a23-937f-148f4107a8f3','Skywalker','2','3dcbd27e-702e-4673-a4c9-0d0104a34b5b');
	
INSERT INTO public.satellites(id, name, "order", position_id)
	VALUES ('3993b12c-00f7-4ded-9323-002155dbfc5a','Kenobi','1','a1177205-195d-4f32-9b9b-59bf98b84134');
	
INSERT INTO public.satellites(id, name, "order", position_id)
	VALUES ('d917e2e2-2e4f-4079-9727-baf707dd336b','Sato','3','ef460c9e-0fb0-460a-9c3a-15f037158666');

INSERT INTO public.positions(id, x, y)
	VALUES ('ef460c9e-0fb0-460a-9c3a-15f037158666', 500, 100);
	
INSERT INTO public.positions(id, x, y)
	VALUES ('3dcbd27e-702e-4673-a4c9-0d0104a34b5b', 100, -100);
	
INSERT INTO public.positions(id, x, y)
	VALUES ('a1177205-195d-4f32-9b9b-59bf98b84134', -500, -200);