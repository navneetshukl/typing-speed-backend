--
-- PostgreSQL database dump
--

\restrict s65xKDa4LFYp0dneRKRrrnyp5ajeqyz7chgYl6fxGvBkXwzqiNjhGwY5dGF34Eo

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

--
-- Name: set_created_at(); Type: FUNCTION; Schema: public; Owner: navneetshukla
--

CREATE FUNCTION public.set_created_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  IF NEW.created_at IS NULL THEN
    NEW.created_at := CURRENT_TIMESTAMP;
  END IF;
  RETURN NEW;
END;
$$;


ALTER FUNCTION public.set_created_at() OWNER TO navneetshukla;

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
-- Name: users; Type: TABLE; Schema: public; Owner: navneetshukla
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO navneetshukla;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: navneetshukla
--

CREATE SEQUENCE public.users_id_seq
    AS integer
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
-- Name: usertypingdata; Type: TABLE; Schema: public; Owner: navneetshukla
--

CREATE TABLE public.usertypingdata (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id character varying(100) NOT NULL,
    total_error integer NOT NULL,
    total_words integer NOT NULL,
    typed_words integer NOT NULL,
    total_time integer NOT NULL,
    total_time_taken_by_user integer NOT NULL,
    wpm integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.usertypingdata OWNER TO navneetshukla;

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
-- Name: usertypingdata usertypingdata_pkey; Type: CONSTRAINT; Schema: public; Owner: navneetshukla
--

ALTER TABLE ONLY public.usertypingdata
    ADD CONSTRAINT usertypingdata_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: navneetshukla
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: usertypingdata set_created_at_trigger; Type: TRIGGER; Schema: public; Owner: navneetshukla
--

CREATE TRIGGER set_created_at_trigger BEFORE INSERT ON public.usertypingdata FOR EACH ROW EXECUTE FUNCTION public.set_created_at();


--
-- PostgreSQL database dump complete
--

\unrestrict s65xKDa4LFYp0dneRKRrrnyp5ajeqyz7chgYl6fxGvBkXwzqiNjhGwY5dGF34Eo

