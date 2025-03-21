CREATE TABLE public.orders (
    id bigint NOT NULL,
    user_id uuid NOT NULL,
    status varchar(255) NOT NULL,
    accrual numeric DEFAULT 0 NOT NULL,
    uploaded_at timestamp DEFAULT now() NOT NULL,
    CONSTRAINT orders_pk PRIMARY KEY (id),
    CONSTRAINT orders_users_fk FOREIGN KEY (user_id) REFERENCES public.users (id)
);
CREATE INDEX orders_user_id_idx ON public.orders USING btree (user_id);
