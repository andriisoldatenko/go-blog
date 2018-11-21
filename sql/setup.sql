create table users (
  id         serial primary key,
  name       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null
);

create table posts (
  id         serial primary key,
  body
  user_id
  thread_id  integer references threads(id),
  created_at timestamp not null
  varchar(64) not null unique,
  text,
  integer references users(id),
);