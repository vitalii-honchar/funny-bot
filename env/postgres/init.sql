create table tg_user (
    "id" serial primary key,
    "first_name" varchar(100) not null,
    "last_name" varchar(100) not null,
    "username" varchar(100) not null,
    "chat_id" bigint not null unique,
    "notification_time" timestamp not null
);