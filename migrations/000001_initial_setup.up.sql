create table author (
    author_id UUID not null primary key,
    name varchar(48) not null
);

create table category (
    category_id uuid not null primary key,
    name varchar(128) not null,
    parent_uuid uuid
);

create table book (
    book_id UUID not null primary key,
    name varchar(128) not null,
    author_id uuid not null references author(author_id),
    parent_uuid uuid not null references category(parent_uuid)
);