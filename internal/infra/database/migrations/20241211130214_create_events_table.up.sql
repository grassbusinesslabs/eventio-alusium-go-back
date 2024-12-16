CREATE TABLE IF NOT EXISTS public.events
(
    id              serial PRIMARY KEY,
    user_id         int NOT NULL references public.users (id),
    title           VARCHAR(120) NOT NULL,
    description     text NOT NULL,
    status          VARCHAR(30) NOT NULL,
    date            timestamptz NOT NULL,
    image           VARCHAR(120) NOT NULL,
    location        VARCHAR(120) NOT NULL,
    lat             float NOT NULL,
    lon             float NOT NULL,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz NULL
)
