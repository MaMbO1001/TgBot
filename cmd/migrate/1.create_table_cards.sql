-- up
create table users_table
(
    id         uuid   not null default gen_random_uuid(),
    name       text   not null,
    tg_user_id BIGINT not null unique,

    primary key (tg_user_id)
);

create table folder
(
    id         uuid                          not null default gen_random_uuid(),
    name       text                          not null unique,
    tg_user_id BIGINT references users_table not null,

    primary key (id)
);

create table card
(
    id        uuid                   not null default gen_random_uuid(),
    name      text                   not null,
    text      text                   not null default '',
    folder_id uuid references folder not null,

    primary key (id)
);

-- down

drop table card;
drop table folder;
drop table users_table;