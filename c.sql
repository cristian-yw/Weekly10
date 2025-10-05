--
-- PostgreSQL database dump
--

\restrict ievNEj4X7f3kAN2btTsupYvyslXrBtamnGjSOUpTATSdK4BashaJJuzekZ4QW8C

-- Dumped from database version 18.0
-- Dumped by pg_dump version 18.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: order_status; Type: TYPE; Schema: public; Owner: cegans
--

CREATE TYPE public.order_status AS ENUM (
    'paid',
    'pending',
    'cancel'
);


ALTER TYPE public.order_status OWNER TO cegans;

--
-- Name: type_role; Type: TYPE; Schema: public; Owner: cegans
--

CREATE TYPE public.type_role AS ENUM (
    'user',
    'admin'
);


ALTER TYPE public.type_role OWNER TO cegans;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name character varying(50)
);


ALTER TABLE public.categories OWNER TO cegans;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO cegans;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: cinemas; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.cinemas (
    id integer NOT NULL,
    name character varying(255),
    image_url character varying(255)
);


ALTER TABLE public.cinemas OWNER TO cegans;

--
-- Name: cinemas_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.cinemas_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cinemas_id_seq OWNER TO cegans;

--
-- Name: cinemas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.cinemas_id_seq OWNED BY public.cinemas.id;


--
-- Name: genres; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.genres (
    id integer NOT NULL,
    tmdb_id integer,
    name character varying(100)
);


ALTER TABLE public.genres OWNER TO cegans;

--
-- Name: genres_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.genres_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.genres_id_seq OWNER TO cegans;

--
-- Name: genres_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.genres_id_seq OWNED BY public.genres.id;


--
-- Name: locations; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.locations (
    id integer NOT NULL,
    location character varying(100)
);


ALTER TABLE public.locations OWNER TO cegans;

--
-- Name: locations_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.locations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.locations_id_seq OWNER TO cegans;

--
-- Name: locations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.locations_id_seq OWNED BY public.locations.id;


--
-- Name: movies; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.movies (
    id integer NOT NULL,
    tmdb_id integer,
    title character varying(255),
    overview text,
    release_date date,
    runtime integer,
    poster_path character varying(255),
    backdrop_path character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    popularity numeric,
    vote_average numeric,
    vote_count integer
);


ALTER TABLE public.movies OWNER TO cegans;

--
-- Name: movies_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.movies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.movies_id_seq OWNER TO cegans;

--
-- Name: movies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.movies_id_seq OWNED BY public.movies.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    user_id integer,
    schedule_id integer,
    order_date timestamp without time zone,
    total_price numeric(10,2),
    status public.order_status
);


ALTER TABLE public.orders OWNER TO cegans;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO cegans;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: payments; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.payments (
    id integer NOT NULL,
    order_id integer,
    payment_method_id integer,
    virtual_account character varying(20),
    amount numeric(10,2),
    status public.order_status,
    paid_at timestamp without time zone,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


ALTER TABLE public.payments OWNER TO cegans;

--
-- Name: payments_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.payments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payments_id_seq OWNER TO cegans;

--
-- Name: payments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.payments_id_seq OWNED BY public.payments.id;


--
-- Name: payments_method; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.payments_method (
    id integer NOT NULL,
    name character varying(50)
);


ALTER TABLE public.payments_method OWNER TO cegans;

--
-- Name: payments_method_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.payments_method_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payments_method_id_seq OWNER TO cegans;

--
-- Name: payments_method_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.payments_method_id_seq OWNED BY public.payments_method.id;


--
-- Name: persons; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.persons (
    id integer NOT NULL,
    tmdb_id integer,
    name character varying(100),
    profile_path character varying(255)
);


ALTER TABLE public.persons OWNER TO cegans;

--
-- Name: persons_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.persons_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.persons_id_seq OWNER TO cegans;

--
-- Name: persons_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.persons_id_seq OWNED BY public.persons.id;


--
-- Name: schedules; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.schedules (
    id integer NOT NULL,
    movie_id integer,
    cinema_id integer,
    location_id integer,
    time_id integer,
    date date
);


ALTER TABLE public.schedules OWNER TO cegans;

--
-- Name: schedules_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.schedules_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.schedules_id_seq OWNER TO cegans;

--
-- Name: schedules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.schedules_id_seq OWNED BY public.schedules.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO cegans;

--
-- Name: times; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.times (
    id integer NOT NULL,
    start_time timestamp without time zone
);


ALTER TABLE public.times OWNER TO cegans;

--
-- Name: times_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.times_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.times_id_seq OWNER TO cegans;

--
-- Name: times_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.times_id_seq OWNED BY public.times.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: cegans
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying(100) NOT NULL,
    password_hash text NOT NULL,
    role public.type_role,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    avatar_url text,
    first_name character varying(100),
    last_name character varying(100),
    phone character varying(20)
);


ALTER TABLE public.users OWNER TO cegans;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: cegans
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO cegans;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cegans
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: cinemas id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.cinemas ALTER COLUMN id SET DEFAULT nextval('public.cinemas_id_seq'::regclass);


--
-- Name: genres id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.genres ALTER COLUMN id SET DEFAULT nextval('public.genres_id_seq'::regclass);


--
-- Name: locations id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.locations ALTER COLUMN id SET DEFAULT nextval('public.locations_id_seq'::regclass);


--
-- Name: movies id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.movies ALTER COLUMN id SET DEFAULT nextval('public.movies_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: payments id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.payments ALTER COLUMN id SET DEFAULT nextval('public.payments_id_seq'::regclass);


--
-- Name: payments_method id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.payments_method ALTER COLUMN id SET DEFAULT nextval('public.payments_method_id_seq'::regclass);


--
-- Name: persons id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.persons ALTER COLUMN id SET DEFAULT nextval('public.persons_id_seq'::regclass);


--
-- Name: schedules id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.schedules ALTER COLUMN id SET DEFAULT nextval('public.schedules_id_seq'::regclass);


--
-- Name: times id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.times ALTER COLUMN id SET DEFAULT nextval('public.times_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.categories (id, name) FROM stdin;
\.


--
-- Data for Name: cinemas; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.cinemas (id, name, image_url) FROM stdin;
\.


--
-- Data for Name: genres; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.genres (id, tmdb_id, name) FROM stdin;
\.


--
-- Data for Name: locations; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.locations (id, location) FROM stdin;
\.


--
-- Data for Name: movies; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.movies (id, tmdb_id, title, overview, release_date, runtime, poster_path, backdrop_path, created_at, updated_at, popularity, vote_average, vote_count) FROM stdin;
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.orders (id, user_id, schedule_id, order_date, total_price, status) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.payments (id, order_id, payment_method_id, virtual_account, amount, status, paid_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: payments_method; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.payments_method (id, name) FROM stdin;
\.


--
-- Data for Name: persons; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.persons (id, tmdb_id, name, profile_path) FROM stdin;
\.


--
-- Data for Name: schedules; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.schedules (id, movie_id, cinema_id, location_id, time_id, date) FROM stdin;
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.schema_migrations (version, dirty) FROM stdin;
12	f
\.


--
-- Data for Name: times; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.times (id, start_time) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: cegans
--

COPY public.users (id, email, password_hash, role, created_at, updated_at, avatar_url, first_name, last_name, phone) FROM stdin;
1	user1@mail.com	$2a$10$er8qgrsfDS/W.SJRHSGg3eXcrZmD3ni9lrTSlyz50pxLIlXmliwri	user	2025-10-05 13:43:44.832423	2025-10-05 13:43:44.832423	\N	\N	\N	\N
\.


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.categories_id_seq', 1, false);


--
-- Name: cinemas_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.cinemas_id_seq', 1, false);


--
-- Name: genres_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.genres_id_seq', 1, false);


--
-- Name: locations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.locations_id_seq', 1, false);


--
-- Name: movies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.movies_id_seq', 1, false);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.orders_id_seq', 1, false);


--
-- Name: payments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.payments_id_seq', 1, false);


--
-- Name: payments_method_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.payments_method_id_seq', 1, false);


--
-- Name: persons_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.persons_id_seq', 1, false);


--
-- Name: schedules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.schedules_id_seq', 1, false);


--
-- Name: times_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.times_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: cegans
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: cinemas cinemas_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.cinemas
    ADD CONSTRAINT cinemas_pkey PRIMARY KEY (id);


--
-- Name: genres genres_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.genres
    ADD CONSTRAINT genres_pkey PRIMARY KEY (id);


--
-- Name: genres genres_tmdb_id_key; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.genres
    ADD CONSTRAINT genres_tmdb_id_key UNIQUE (tmdb_id);


--
-- Name: locations locations_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.locations
    ADD CONSTRAINT locations_pkey PRIMARY KEY (id);


--
-- Name: movies movies_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_pkey PRIMARY KEY (id);


--
-- Name: movies movies_tmdb_id_key; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_tmdb_id_key UNIQUE (tmdb_id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: payments_method payments_method_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.payments_method
    ADD CONSTRAINT payments_method_pkey PRIMARY KEY (id);


--
-- Name: payments payments_order_id_key; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_order_id_key UNIQUE (order_id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: payments payments_virtual_account_key; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_virtual_account_key UNIQUE (virtual_account);


--
-- Name: persons persons_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.persons
    ADD CONSTRAINT persons_pkey PRIMARY KEY (id);


--
-- Name: persons persons_tmdb_id_key; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.persons
    ADD CONSTRAINT persons_tmdb_id_key UNIQUE (tmdb_id);


--
-- Name: schedules schedules_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.schedules
    ADD CONSTRAINT schedules_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: times times_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.times
    ADD CONSTRAINT times_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: orders orders_schedule_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES public.schedules(id);


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: cegans
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict ievNEj4X7f3kAN2btTsupYvyslXrBtamnGjSOUpTATSdK4BashaJJuzekZ4QW8C

