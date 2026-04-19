--
-- PostgreSQL database dump
--

\restrict 1SMfukc4yRBpSVOf6YU6cfZIzujuRCCFh23RPwdZWmiLiwn5PMrcr6etVoSTSPg

-- Dumped from database version 16.13
-- Dumped by pg_dump version 16.13

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

ALTER TABLE IF EXISTS ONLY public.user_rewards DROP CONSTRAINT IF EXISTS user_rewards_quest_id_fkey;
ALTER TABLE IF EXISTS ONLY public.user_quest_progress DROP CONSTRAINT IF EXISTS user_quest_progress_quest_id_fkey;
ALTER TABLE IF EXISTS ONLY public.quest_steps DROP CONSTRAINT IF EXISTS quest_steps_quest_id_fkey;
DROP INDEX IF EXISTS public.idx_users_email;
DROP INDEX IF EXISTS public.idx_user_progress_user;
DROP INDEX IF EXISTS public.idx_quests_slug;
DROP INDEX IF EXISTS public.idx_quest_steps_quest;
DROP INDEX IF EXISTS public.idx_folks_name;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS ONLY public.user_rewards DROP CONSTRAINT IF EXISTS user_rewards_pkey;
ALTER TABLE IF EXISTS ONLY public.user_quest_progress DROP CONSTRAINT IF EXISTS user_quest_progress_user_id_quest_id_key;
ALTER TABLE IF EXISTS ONLY public.user_quest_progress DROP CONSTRAINT IF EXISTS user_quest_progress_pkey;
ALTER TABLE IF EXISTS ONLY public.schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.quests DROP CONSTRAINT IF EXISTS quests_slug_key;
ALTER TABLE IF EXISTS ONLY public.quests DROP CONSTRAINT IF EXISTS quests_pkey;
ALTER TABLE IF EXISTS ONLY public.quest_steps DROP CONSTRAINT IF EXISTS quest_steps_quest_id_step_id_key;
ALTER TABLE IF EXISTS ONLY public.quest_steps DROP CONSTRAINT IF EXISTS quest_steps_pkey;
ALTER TABLE IF EXISTS ONLY public.folks DROP CONSTRAINT IF EXISTS folks_pkey;
ALTER TABLE IF EXISTS ONLY public.folks DROP CONSTRAINT IF EXISTS folks_name_key;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.user_rewards;
DROP TABLE IF EXISTS public.user_quest_progress;
DROP TABLE IF EXISTS public.schema_migrations;
DROP TABLE IF EXISTS public.quests;
DROP TABLE IF EXISTS public.quest_steps;
DROP TABLE IF EXISTS public.folks;
DROP EXTENSION IF EXISTS "uuid-ossp";
--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: folks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.folks (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    lat numeric(9,6),
    lon numeric(9,6),
    title character varying(120),
    summary text,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: quest_steps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.quest_steps (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    quest_id uuid NOT NULL,
    step_id character varying(50) NOT NULL,
    step_order integer NOT NULL,
    step_type character varying(30) NOT NULL,
    title character varying(255) NOT NULL,
    content jsonb NOT NULL,
    on_success jsonb
);


--
-- Name: quests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.quests (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    slug character varying(100) NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    cover_url character varying(500),
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: user_quest_progress; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_quest_progress (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    quest_id uuid NOT NULL,
    current_step_id character varying(50),
    completed_steps jsonb DEFAULT '[]'::jsonb,
    status character varying(20) DEFAULT 'in_progress'::character varying,
    started_at timestamp with time zone DEFAULT now(),
    completed_at timestamp with time zone
);


--
-- Name: user_rewards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_rewards (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    quest_id uuid,
    reward_type character varying(30) NOT NULL,
    reward_key character varying(100) NOT NULL,
    granted_at timestamp with time zone DEFAULT now(),
    metadata jsonb
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    username character varying(100),
    role character varying(50) DEFAULT 'user'::character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Data for Name: folks; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.folks (id, name, lat, lon, title, summary, created_at) FROM stdin;
9e13c30f-e605-44f7-9ee6-5c1a8e749663	Чукотский автономный округ	64.734200	177.518900	Эскимосы	{\r\n  "text": "**Древняя художественная традиция**\\n\\nНаходки археологов подтверждают существование изобразительного искусства на побережье Чукотки с древнейших времён.\\n  - **сложные узоры** на амулетах, украшениях, наконечниках гарпунов\\n  - **орнамент** с сакральным значением\\n  - **резьба по кости**: фигурки людей и животных, сочетающие достоверность и фантазию\\n  - косторезами становились **потомственные зверобои**\\n\\n**Современное декоративно-прикладное искусство**\\n\\nИзделия мастеров Чукотки практичны и выразительны:\\n  - **шкатулки** и **футляры для очков**\\n  - **чехлы для мобильных телефонов**\\n  - другие предметы быта с традиционным орнаментом",\r\n  "images": []\r\n}	2026-04-17 20:13:01.218924+00
93e5afed-bde9-4047-98f4-7172faf91565	Красноярский край	56.015280	92.893250	Энцы	{\r\n  "text": "**Традиционные ремёсла**\\n\\nС незапамятных времён энцы владели технологиями выделки шкуры:\\n  - **сыромятную кожу**\\n  - **дубленую кожу**\\n  - **ровдугу** (замша из оленьей шкуры)\\n\\n**Декоративное искусство**\\n\\nОрнамент энцев очень разнообразен. Наиболее распространены:\\n  - узоры из кусочков **оленьей шерсти** тёмного и светлого цвета\\n  - композиции из **разноцветного сукна**",\r\n  "images": []\r\n}	2026-04-17 19:50:16.08323+00
aa851270-06c3-49e3-b467-cb63ba52a02a	Ленинградская область	59.576400	30.123900	Водь	{\r\n  "text": "**Водь** — коренной малочисленный финно-угорский народ, проживающий в Ленинградской области.\\n\\n**Традиционный костюм**\\n\\nВожанки славились умением изготавливать домотканую одежду. Для традиционного наряда характерны:\\n  - **льняные передники**\\n  - **суконные передники**\\n\\n**Декоративные украшения**\\n\\nКостюм дополняли сложные нагрудные украшения (риссико, мюэтси), декорированные:\\n  - **вышивкой** и **тесьмой**\\n  - **бисером**\\n  - **раковинами каури**",\r\n  "images": []\r\n}	2026-04-17 20:29:46.29735+00
368f0cdd-d671-4d7a-89cf-bb88ec20f62a	Приморский край	43.105600	131.873500	Удэгейцы	{\r\n  "text": "**Традиционные ремёсла удэгейцев**\\n\\nРемёсла этого коренного народа Дальнего Востока тесно связаны с лесным хозяйством и охотой.\\n\\n**Ключевые промыслы**\\n\\n  - **обработка кожи и меха** (изготовление одежды и обуви)\\n  - **резьба по дереву и бересте** (создание посуды и оморочек)\\n  - **художественное оформление бытовых предметов**",\r\n  "images": []\r\n}	2026-04-17 20:26:03.587547+00
862e59ec-f61e-4df8-a6a3-c715a898ee67	Челябинская область	54.358333	60.501667	Нагайбаки	{\r\n  "text": "**Песенный фольклор**\\n\\nПесни нагайбаков отражали повседневные заботы и стали важной частью быта. В каждом селе сохраняются свои уникальные напевы.\\n\\n**Традиционная кухня**\\n\\nОснова рациона — молочные, мясные и зерновые продукты. Характерны:\\n  - **выпечка** из пшеничной муки\\n  - **болтушка** и мучные супы\\n  - свежие овощи, фрукты, **компоты** и **кисели**\\n\\n**Хозяйство и промыслы**\\n\\nОснова экономики — степное земледелие (пшеница, овёс, ячмень) и животноводство. Также развиты:\\n  - **пчеловодство**\\n  - **плотницкое дело**\\n  - **производство экипажей**",\r\n  "images": []\r\n}	2026-04-18 21:21:14.552229+00
\.


--
-- Data for Name: quest_steps; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.quest_steps (id, quest_id, step_id, step_order, step_type, title, content, on_success) FROM stdin;
ac7b0ed6-34fa-4db0-abdc-03f71b69a746	736bba5d-aa9f-4f61-b915-f8d5ce891751	step_1	1	intro	Знак рода	{"text": "У каждого рода хантов есть свой знак — тамга. Её вырезали на оленьих ушах, наносили на утварь и оружие. Тамга оберегала и помогала узнать свою вещь среди сотен других.", "image": "/static/assets/khanty/tamga_example.png"}	\N
615396b9-f6dd-4309-9f60-c03119b1e586	736bba5d-aa9f-4f61-b915-f8d5ce891751	step_2	2	quiz	Проверь знания	{"options": [{"id": "a", "text": "Тепло и жизнь", "correct": true}, {"id": "b", "text": "Охотничья удача", "correct": false}, {"id": "c", "text": "Защита от духов", "correct": false}], "question": "Что означает узор «солнце» в орнаменте хантов?"}	{"reward": {"badge": "novice"}, "unlock": ["step_3"]}
9c7674a9-f5c4-48dd-805f-4c88c5a40133	7224de64-6966-4425-8840-43d8883a6fff	step_1	1	intro	Знакомство с игрушками	{"text": "В России каждый регион славится своими уникальными игрушками. Сегодня ты познакомишься с четырьмя традиционными промыслами и узнаешь, откуда они родом."}	\N
79ab7025-d972-49ec-9b97-21d222959da1	736bba5d-aa9f-4f61-b915-f8d5ce891751	step_3	3	builder	Создай узор	{"base": "empty", "goal": "Собери традиционный хантыйский узор. Обязательно добавь солнце!", "patterns": [{"id": "solnce", "file": "/static/assets/khanty/solnce.jpg", "name": "Солнце", "symbol": "☀️"}, {"id": "lebed", "file": "/static/assets/khanty/lebed.jpg", "name": "Лебедь", "symbol": "🦢"}, {"id": "dom", "file": "/static/assets/khanty/dom.jpg", "name": "Дом", "symbol": "🏠"}, {"id": "shishka", "file": "/static/assets/khanty/shishka.jpg", "name": "Шишка", "symbol": "🌲"}], "required": ["solnce"]}	{"reward": {"badge": "craftsman", "points": 100}}
75ba83f3-b010-4eb9-b141-6101912dd832	7224de64-6966-4425-8840-43d8883a6fff	step_6	6	map_match	Помоги игрушкам вернуться в родной край	{"goal": "Нажми на игрушку, а затем на правильный регион на карте.", "toys": [{"id": "tulskaya", "name": "Тульская игрушка", "image": "/static/assets/toys/tulskaya.jpg", "description": "Глиняная расписная игрушка с яркими полосками и узорами."}, {"id": "mazykskaya", "name": "Мазыкская игрушка", "image": "/static/assets/toys/mazykskaya.jpg", "description": "Яркая глиняная игрушка."}, {"id": "vyrkovskaya", "name": "Вырковская игрушка", "image": "/static/assets/toys/vyrkovskaya.jpg", "description": "Глиняные коньки и птицы."}, {"id": "romanovskaya", "name": "Романовская игрушка", "image": "/static/assets/toys/romanovskaya.jpg", "description": "Глиняная игрушка с характерными узорами."}], "regions": [{"id": "tula_region", "name": "Тульская область", "center": [54.2, 37.6], "correct_toys": ["tulskaya"]}, {"id": "vladimir_region", "name": "Владимирская область", "center": [56.1, 40.4], "correct_toys": ["mazykskaya"]}, {"id": "ryazan_region", "name": "Рязанская область", "center": [54.6, 39.7], "correct_toys": ["vyrkovskaya"]}, {"id": "lipetsk_region", "name": "Липецкая область", "center": [52.6, 39.6], "correct_toys": ["romanovskaya"]}]}	\N
d3846b0f-5113-478e-9df2-961a3fd8dc35	7224de64-6966-4425-8840-43d8883a6fff	step_3	3	intro	Мазыкская игрушка	{"text": "<img src='/static/assets/toys/mazykskaya.jpg' alt='Мазыкская игрушка' style='max-width: 100%; border-radius: 8px; margin: 16px 0; display: block;'><p><h3>Минималистичные коники</h3><strong>Промысел владимирских мастеров</strong>, названный в честь таинственного сообщества мазыков (офеней) — потомков скоморохов со своим языком и тайными знаниями.</p><p>Главная особенность: игрушку вырезают <strong>только топориком</strong>! Мастер так натачивал инструмент, что мог создавать тончайшие детали.</p><p>Мастера не скрывали дефекты дерева, а <strong>превращали сучки в хвосты, клювы и конечности</strong> фигурок. Игрушки получались грубоватые, лаконичные, угловатые.</p><p><strong>Фигурки не расписывают</strong> — оставляют натуральное дерево. Самый частый персонаж — <em>коник</em> (лошадь).</p><p>Игрушки делали <strong>не на продажу, а только для своих детей</strong>. Сейчас промысел почти исчез — им занимаются лишь несколько мастеров.</p>"}	\N
759800fb-d3e7-4d0c-8a3d-ed11c68867ec	7224de64-6966-4425-8840-43d8883a6fff	step_2	2	intro	Тульская городская игрушка	{"text": "<img src='/static/assets/toys/tulskaya.jpg' alt='Тульская городская игрушка' style='max-width: 100%; border-radius: 8px; margin: 16px 0; display: block;'><h3>Барыни под зонтиками</h3><p>Основной сюжет тульской игрушки можно охарактеризовать как <strong>праздный городской мир глазами кустарного мастера</strong>.</p><p>Среди героев нехарактерные для кустарной игрушки городские персонажи — разодетые барыни, няни, кормилицы, военные, дамы и монахи, получившие общее ироничное название <em>«князьки»</em>.</p><p>Ещё одна особенность игрушек — <strong>вытянутая, «долговязая» форма</strong>. Фигурки традиционно фронтальные, но проработаны со всех сторон.</p>"}	\N
ac51401a-a386-42f5-8184-f5324c2a552f	7224de64-6966-4425-8840-43d8883a6fff	step_4	4	intro	Вырковская игрушка	{"text": "<img src='/static/assets/toys/vyrkovskaya.jpg' alt='Вырковская игрушка' style='max-width: 100%; border-radius: 8px; margin: 16px 0; display: block;'><h3>Расскажут о Рязани</h3><p><strong>Глиняные свистульки из ярко-красной глины</strong>, которую добывали рядом с деревней Вырково. В начале XX века это стало <strong>семейным делом</strong> — детей привлекали к производству с 11–12 лет.</p><p>Мастера делали коньков, петушков, медведей. Особо искусные создавали <strong>«говорящих» зверей</strong>: свинки хрюкали, собачки лаяли, а птицы чирикали!</p><p>Игрушка <strong>минималистична</strong>: простые наивные формы и <strong>отсутствие ярких красок</strong> — только натуральная красная глина.</p><p>Увидеть изделия можно в Касимовском краеведческом и Рязанском художественном музеях.</p>"}	\N
6307c81b-7c8a-4b06-8712-9c76707fec18	7224de64-6966-4425-8840-43d8883a6fff	step_5	5	intro	Романовская игрушка	{"text": "<img src='/static/assets/toys/romanovskaya.jpg' alt='Романовская игрушка' style='max-width: 100%; border-radius: 8px; margin: 16px 0; display: block;'><h3>Ах какие глазки!</h3><p><strong>Одна из старейших глиняных игрушек России</strong> — первые образцы датируются <strong>XVII веком</strong>. В отличие от других регионов, её сразу выделили в отдельный промысел и продавали на ярмарках. Ласково называют <em>«романушка»</em>.</p><p>Более <strong>50 сюжетов</strong>: от животных до исторических фигур (Николай II, Лев Толстой). Почти всегда это <strong>свистулька с отверстиями</strong> — можно воспроизводить сложные мелодии благодаря особой акустике местной глины.</p><p>Легко узнать по <strong>схематичному лицу с отверстиями</strong> вместо глаз, носа и рта. Традиционно покрывали <strong>тёмно-зелёной или коричневой глазурью</strong>, позже стали расписывать яркими красками и серебрить детали.</p><p>Промысел <strong>живёт до сих пор</strong> в селе Троицкое Липецкой области!</p>"}	\N
\.


--
-- Data for Name: quests; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.quests (id, slug, title, description, cover_url, is_active, created_at, updated_at) FROM stdin;
736bba5d-aa9f-4f61-b915-f8d5ce891751	khanty-crafts	Тропа мастеров Ханты	Познакомься с древними ремёслами и традициями хантов. Узнай о тамгах, узорах и создай свой оберег.	/static/assets/khanty/cover.jpg	t	2026-04-17 22:18:50.199766+00	2026-04-17 22:18:50.199766+00
7224de64-6966-4425-8840-43d8883a6fff	toy-map-quest	Карта народных игрушек	Распредели традиционные игрушки по их родным регионам.	/static/assets/toys/cover.jpg	t	2026-04-18 11:38:00.499688+00	2026-04-18 11:38:00.499688+00
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.schema_migrations (version, dirty) FROM stdin;
3	f
\.


--
-- Data for Name: user_quest_progress; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.user_quest_progress (id, user_id, quest_id, current_step_id, completed_steps, status, started_at, completed_at) FROM stdin;
\.


--
-- Data for Name: user_rewards; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.user_rewards (id, user_id, quest_id, reward_type, reward_key, granted_at, metadata) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, email, password_hash, username, role, created_at, updated_at) FROM stdin;
a8cb7f26-be48-4e1a-9192-10b1f45fd0fc	example@mail.ru	$2a$10$eYKg9WZpc7JSf0o.XOVYYucbeIxUsEW3YBMcokW/NcDc753lRp3I2		user	2026-04-17 21:10:38.279853+00	2026-04-17 21:10:38.279853+00
2403991b-1363-43c5-93c6-de648e82e3df	example1@mail.ru	$2a$10$YXhnxfrv208mfL.GnPDVbuUVIaF4NFxJ8TrUSmeMFh1DHWb2Bf8fC		user	2026-04-18 09:23:08.063708+00	2026-04-18 09:23:08.063708+00
d2eaed03-62db-4766-8ee5-41d493f4d3cf	12345@mail.ru	$2a$10$yEDFnUhKiptF7qk7Pihb4uxtoqw9NrJfnigE3M9/Hr0b.9fT/25t2		user	2026-04-18 09:31:23.871397+00	2026-04-18 09:31:23.871397+00
\.


--
-- Name: folks folks_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.folks
    ADD CONSTRAINT folks_name_key UNIQUE (name);


--
-- Name: folks folks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.folks
    ADD CONSTRAINT folks_pkey PRIMARY KEY (id);


--
-- Name: quest_steps quest_steps_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quest_steps
    ADD CONSTRAINT quest_steps_pkey PRIMARY KEY (id);


--
-- Name: quest_steps quest_steps_quest_id_step_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quest_steps
    ADD CONSTRAINT quest_steps_quest_id_step_id_key UNIQUE (quest_id, step_id);


--
-- Name: quests quests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quests
    ADD CONSTRAINT quests_pkey PRIMARY KEY (id);


--
-- Name: quests quests_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quests
    ADD CONSTRAINT quests_slug_key UNIQUE (slug);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: user_quest_progress user_quest_progress_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_quest_progress
    ADD CONSTRAINT user_quest_progress_pkey PRIMARY KEY (id);


--
-- Name: user_quest_progress user_quest_progress_user_id_quest_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_quest_progress
    ADD CONSTRAINT user_quest_progress_user_id_quest_id_key UNIQUE (user_id, quest_id);


--
-- Name: user_rewards user_rewards_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_rewards
    ADD CONSTRAINT user_rewards_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_folks_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_folks_name ON public.folks USING btree (name);


--
-- Name: idx_quest_steps_quest; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_quest_steps_quest ON public.quest_steps USING btree (quest_id, step_order);


--
-- Name: idx_quests_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_quests_slug ON public.quests USING btree (slug);


--
-- Name: idx_user_progress_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_progress_user ON public.user_quest_progress USING btree (user_id, status);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: quest_steps quest_steps_quest_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quest_steps
    ADD CONSTRAINT quest_steps_quest_id_fkey FOREIGN KEY (quest_id) REFERENCES public.quests(id) ON DELETE CASCADE;


--
-- Name: user_quest_progress user_quest_progress_quest_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_quest_progress
    ADD CONSTRAINT user_quest_progress_quest_id_fkey FOREIGN KEY (quest_id) REFERENCES public.quests(id) ON DELETE CASCADE;


--
-- Name: user_rewards user_rewards_quest_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_rewards
    ADD CONSTRAINT user_rewards_quest_id_fkey FOREIGN KEY (quest_id) REFERENCES public.quests(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

\unrestrict 1SMfukc4yRBpSVOf6YU6cfZIzujuRCCFh23RPwdZWmiLiwn5PMrcr6etVoSTSPg

