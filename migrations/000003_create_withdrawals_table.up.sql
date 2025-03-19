CREATE TABLE public.withdrawals (
    id bigserial NOT NULL,
    user_id uuid NOT NULL,
    order_id int8 NOT NULL,
    sum numeric NOT NULL,
    processed_at timestamp DEFAULT now() NOT NULL,
    CONSTRAINT withdrawals_pk PRIMARY KEY (id),
    CONSTRAINT withdrawals_users_fk FOREIGN KEY (user_id) REFERENCES public.users (id)
    -- CONSTRAINT withdrawals_orders_fk FOREIGN KEY (order_id) REFERENCES public.orders (id)
);
CREATE INDEX withdrawals_user_id_idx ON public.withdrawals USING btree (user_id);
CREATE INDEX withdrawals_order_id_idx ON public.withdrawals USING btree (order_id);
