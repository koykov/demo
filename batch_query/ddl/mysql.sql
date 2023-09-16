drop table if exists users;

create table users
(
    id      int auto_increment primary key,
    name    varchar(128) null,
    status  int          null,
    bio     text         null,
    balance float        null
);
