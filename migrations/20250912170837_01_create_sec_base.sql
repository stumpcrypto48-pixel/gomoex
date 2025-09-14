-- Add migration script here


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists moex.secs (
    id uuid default uuid_generate_v4() primary key,
    sec_id varchar(256) not null,
    short_name varchar(256),
    regnumber varchar(256),
    name varchar(256),
    isin varchar(256),
    is_traded int,
    emitent_id int,
    emitent_title varchar(256),
    emitent_inn varchar(256),
    emitent_okpo varchar(256),
    type varchar(256),
    group_name varchar(256),
    primary_board_id varchar(256),
    market_price_board_id varchar(256),
    create_date timestamp with time zone default now(),
    update_date timestamp with time zone default now(),
    delete_date timestamp with time zone default null
);
