drop table posts;
drop table threads;
drop table sessions;
drop table users;


create table users (
  id         integer primary key,
  uuid       text not null unique,
  name       text,
  email      text not null unique,
  password   text not null,
  created_at timestamp not null   
);

create table sessions (
  id         integer primary key,
  uuid       text not null unique,
  email      text,
  user_id    integer references users(id),
  created_at text not null   
);

create table threads (
  id         integer primary key,
  uuid       text not null unique,
  topic      text,
  user_id    integer references users(id),
  created_at text not null       
);

create table posts (
  id         integer primary key,
  uuid       text not null unique,
  body       text,
  user_id    integer references users(id),
  thread_id  integer references threads(id),
  created_at text not null  
);