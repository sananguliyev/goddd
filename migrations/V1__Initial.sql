--
-- Name: ratings; Type: TABLE; Schema: public; Owner: goddd
--

CREATE TABLE IF NOT EXISTS ratings (
                         id uuid NOT NULL,
                         recipe_id uuid NOT NULL,
                         value smallint NOT NULL,
                         PRIMARY KEY (id)
);


ALTER TABLE ratings OWNER TO goddd;

--
-- Name: src; Type: TABLE; Schema: public; Owner: goddd
--

CREATE TABLE IF NOT EXISTS recipes (
                         id uuid NOT NULL,
                         name character varying NOT NULL,
                         prepare_time smallint NOT NULL,
                         difficulty smallint NOT NULL,
                         is_vegetarian boolean NOT NULL,
                         PRIMARY KEY (id)
);


ALTER TABLE recipes OWNER TO goddd;

--
-- Data for Name: ratings; Type: TABLE DATA; Schema: public; Owner: goddd
--

COPY ratings (id, recipe_id, value) FROM stdin;
341f1ac7-7092-423d-849b-370a665b9769	0d93bda0-f040-4564-967b-e59bf5571dcd	4
5fd69e74-48b6-447f-850c-05bbf5f15717	0d93bda0-f040-4564-967b-e59bf5571dcd	5
902a679e-9815-4573-a8d0-d3ecae652006	0d93bda0-f040-4564-967b-e59bf5571dcd	5
\.


--
-- Data for Name: src; Type: TABLE DATA; Schema: public; Owner: goddd
--

COPY recipes (name, prepare_time, difficulty, is_vegetarian, id) FROM stdin;
Dolma	45	5	f	d03bb3ce-23f1-40e4-9a9e-59d9dadac472
Ash	45	2	t	455b8d23-0860-4e7f-ae5e-30227c603419
Piti	60	4	f	b399cdc6-fb8c-49ca-8dc3-fca9f9dbf62b
Lyulya kebab	15	3	f	61b4e6c2-8d3d-4664-8588-f8554b9504f5
Erishte	45	4	t	81ba1608-8920-4c65-ab52-838896b12383
Dovgha	20	2	t	0d93bda0-f040-4564-967b-e59bf5571dcd
\.
