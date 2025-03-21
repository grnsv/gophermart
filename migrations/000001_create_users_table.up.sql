CREATE TABLE public.users (
    id uuid NOT NULL,
    login varchar(255) NOT NULL,
    "password" varchar(255) NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_unique UNIQUE (login)
);
