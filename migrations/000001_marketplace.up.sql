create extension if not exists "uuid-ossp";

create table users (
    id uuid primary key default uuid_generate_v4(),
    login varchar(50) not null unique,
    password text not null,
    created_at timestamp not null default now()
);

create table advertisements (
    id uuid primary key default uuid_generate_v4(),
    title text not null,
    content text not null,
    image_url text,
    price numeric(11, 2) not null,
    user_id uuid not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    foreign key (user_id) references users (id) on delete cascade
);
