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


create table tasks
(
    task_id serial
        constraint tasks_pk
            primary key,
    public_id uuid,
    author uuid,
    title varchar(255),
    description text,
    assign_cost int,
    done_cost int,
    status varchar(10),
    created_at timestamp,
    assigned_to uuid,
    assigned_at timestamp
);

create index tasks_author_index
    on tasks (author);

create index tasks_public_id_index
    on tasks (public_id);

create index tasks_created_at_index
    on tasks (created_at);

create index tasks_assigned_to_index
    on tasks (assigned_to);
