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
);

create table book_category (
    id uuid not null primary key,
    book_id uuid not null references book(book_id),
    category_id uuid not null references category(category_id)
);