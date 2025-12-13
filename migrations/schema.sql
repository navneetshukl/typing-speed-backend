--
-- PostgreSQL database dump
--

\restrict z02vggQp0WzoyT76whqXkJwj2Ap7N5xkE6s3dBKtKewYlD3fKRrr4s0ByZh5dfq

-- Dumped from database version 14.19 (Homebrew)
-- Dumped by pg_dump version 14.19 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: navneetshukla
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO navneetshukla;

--
-- Name: user_typing_data; Type: TABLE; Schema: public; Owner: navneetshukla
--

CREATE TABLE public.user_typing_data (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(100) NOT NULL,
    total_error integer NOT NULL,
    total_words integer NOT NULL,
    typed_words integer NOT NULL,
    total_time integer NOT NULL,
    total_time_taken_by_user integer NOT NULL,
    wpm integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_typing_data OWNER TO navneetshukla;

--
-- Name: users; Type: TABLE; Schema: public; Owner: navneetshukla
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    avg_speed integer DEFAULT 0,
    avg_accuracy integer DEFAULT 0,
    total_test integer DEFAULT 0,
    level integer DEFAULT 0,
    last_test_time timestamp with time zone,
    streak integer DEFAULT 0,
    best_speed integer DEFAULT 0,
    avg_performance integer DEFAULT 0,
    CONSTRAINT users_avg_accuracy_check CHECK (((avg_accuracy >= 0) AND (avg_accuracy <= 100))),
    CONSTRAINT users_avg_speed_check CHECK ((avg_speed >= 0)),
    CONSTRAINT users_level_check CHECK ((level >= 0)),
    CONSTRAINT users_streak_check CHECK ((streak >= 0)),
    CONSTRAINT users_total_test_check CHECK ((total_test >= 0))
);


ALTER TABLE public.users OWNER TO navneetshukla;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: navneetshukla
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO navneetshukla;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: navneetshukla
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: navneetshukla
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: schema_migration schema_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: navneetshukla
--

ALTER TABLE ONLY public.schema_migration
    ADD CONSTRAINT schema_migration_pkey PRIMARY KEY (version);


--
-- Name: user_typing_data user_typing_data_pkey; Type: CONSTRAINT; Schema: public; Owner: navneetshukla
--

ALTER TABLE ONLY public.user_typing_data
    ADD CONSTRAINT user_typing_data_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: navneetshukla
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: navneetshukla
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: navneetshukla
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: user_typing_data fk_user_email; Type: FK CONSTRAINT; Schema: public; Owner: navneetshukla
--

ALTER TABLE ONLY public.user_typing_data
    ADD CONSTRAINT fk_user_email FOREIGN KEY (email) REFERENCES public.users(email) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict z02vggQp0WzoyT76whqXkJwj2Ap7N5xkE6s3dBKtKewYlD3fKRrr4s0ByZh5dfq

