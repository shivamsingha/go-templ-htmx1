CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid(),
    email text not null unique,
    name text not null,
    password text not null,
    email_verified boolean default false,
    created_at timestamp default now(),
    updated_at timestamp default now()
);