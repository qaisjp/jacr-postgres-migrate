--
-- PostgreSQL database dump
--

-- Dumped from database version 9.4.8
-- Dumped by pg_dump version 9.6.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: response_commands; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE response_commands (
    id integer NOT NULL,
    name character varying(32) NOT NULL
);


ALTER TABLE response_commands OWNER TO postgres;

--
-- Name: response_commands_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE response_commands_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE response_commands_id_seq OWNER TO postgres;

--
-- Name: response_commands_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE response_commands_id_seq OWNED BY response_commands.id;


--
-- Name: response_content_commands; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE response_content_commands (
    command integer NOT NULL,
    content integer NOT NULL
);


ALTER TABLE response_content_commands OWNER TO postgres;

--
-- Name: response_contents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE response_contents (
    id integer NOT NULL,
    messages text[] NOT NULL
);


ALTER TABLE response_contents OWNER TO postgres;

--
-- Name: response_contents_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE response_contents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE response_contents_id_seq OWNER TO postgres;

--
-- Name: response_contents_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE response_contents_id_seq OWNED BY response_contents.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE users (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    password character(60) NOT NULL,
    email character varying(254) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    activated boolean DEFAULT false NOT NULL,
    banned boolean DEFAULT false NOT NULL,
    slug character varying(255) NOT NULL,
    level integer DEFAULT 1 NOT NULL
);


ALTER TABLE users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: response_commands id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_commands ALTER COLUMN id SET DEFAULT nextval('response_commands_id_seq'::regclass);


--
-- Name: response_contents id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_contents ALTER COLUMN id SET DEFAULT nextval('response_contents_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Name: response_commands response_commands_id_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_commands
    ADD CONSTRAINT response_commands_id_pkey PRIMARY KEY (id);


--
-- Name: response_commands response_commands_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_commands
    ADD CONSTRAINT response_commands_name_key UNIQUE (name);


--
-- Name: response_content_commands response_content_commands_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_content_commands
    ADD CONSTRAINT response_content_commands_pkey PRIMARY KEY (command, content);


--
-- Name: response_contents response_contents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_contents
    ADD CONSTRAINT response_contents_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_id_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_id_pkey PRIMARY KEY (id);


--
-- Name: users users_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_slug_key UNIQUE (slug);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: response_content_commands response_content_commands_command_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_content_commands
    ADD CONSTRAINT response_content_commands_command_fkey FOREIGN KEY (command) REFERENCES response_commands(id);


--
-- Name: response_content_commands response_content_commands_content_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_content_commands
    ADD CONSTRAINT response_content_commands_content_fkey FOREIGN KEY (content) REFERENCES response_contents(id) ON DELETE CASCADE;


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

