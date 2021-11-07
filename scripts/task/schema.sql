create table users
(
    user_id serial
        constraint users_pk
            primary key,
    public_id uuid,
    user_mail varchar(255),
    user_name varchar,
    user_version smallint,
    user_role varchar(10),
    active bit
);

create index users_active_index
    on users (active);

create index users_public_id_index
    on users (public_id);

create index users_user_role_index
    on users (user_role);

create index users_user_mail_index
    on users (user_mail);