-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION public.set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;
-- +goose StatementEnd

CREATE TABLE public.factions (
    id integer NOT NULL,
    name text,
    description text,
    resources text,
    user_id uuid
);

ALTER TABLE public.factions ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.factions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE public.master_requests (
    id integer NOT NULL,
    to_player integer NOT NULL,
    text_request text DEFAULT ''::text NOT NULL,
    text_response text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    state integer DEFAULT 0 CONSTRAINT master_requests_is_answered_not_null NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER TABLE public.master_requests ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.master_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE public.message_transaction (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    from_chat bigint NOT NULL,
    to_chat bigint[] NOT NULL,
    message_id integer NOT NULL
);

CREATE SEQUENCE public.message_transaction_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.message_transaction_id_seq OWNED BY public.message_transaction.id;

CREATE TABLE public.messages (
    id integer NOT NULL,
    chat_id integer,
    message_title text,
    message_id text,
    message_text text DEFAULT ''::text NOT NULL
);

ALTER TABLE public.messages ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE public.roll_requests (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    title text DEFAULT ''::text CONSTRAINT "roll_requests_Title_not_null" NOT NULL,
    dice_count integer DEFAULT 0 NOT NULL,
    dice_sides integer DEFAULT 0 NOT NULL,
    roll_result integer,
    transaction_id integer NOT NULL
);

ALTER TABLE public.roll_requests ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.roll_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    telegram_name text,
    player_name text DEFAULT ''::text NOT NULL,
    chat_id integer,
    role integer DEFAULT 0 NOT NULL
);

ALTER TABLE ONLY public.message_transaction ALTER COLUMN id SET DEFAULT nextval('public.message_transaction_id_seq'::regclass);

ALTER TABLE ONLY public.factions
    ADD CONSTRAINT factions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.master_requests
    ADD CONSTRAINT master_requests_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.message_transaction
    ADD CONSTRAINT message_transaction_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.roll_requests
    ADD CONSTRAINT roll_requests_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE TRIGGER trigger_set_updated_at BEFORE UPDATE ON public.master_requests FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

ALTER TABLE ONLY public.factions
    ADD CONSTRAINT factions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);

ALTER TABLE ONLY public.roll_requests
    ADD CONSTRAINT fk_master_request FOREIGN KEY (transaction_id) REFERENCES public.master_requests(id) ON DELETE RESTRICT;

ALTER TABLE ONLY public.message_transaction
    ADD CONSTRAINT message_transaction_message_id_fkey FOREIGN KEY (message_id) REFERENCES public.messages(id) ON DELETE CASCADE;

-- +goose Down
BEGIN;

DROP TRIGGER IF EXISTS trigger_set_updated_at ON public.master_requests;

DROP TABLE IF EXISTS public.roll_requests;
DROP TABLE IF EXISTS public.message_transaction;
DROP TABLE IF EXISTS public.factions;

DROP TABLE IF EXISTS public.master_requests;
DROP TABLE IF EXISTS public.messages;
DROP TABLE IF EXISTS public.users;

DROP FUNCTION IF EXISTS public.set_updated_at;

COMMIT;
