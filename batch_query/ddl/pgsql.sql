create schema bq;

create table bq.users
(
    id      serial
        constraint users_pk
        primary key,
    name    varchar(128),
    status  integer,
    bio     text,
    balance double precision
);

