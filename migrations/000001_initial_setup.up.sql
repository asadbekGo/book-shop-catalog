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
    author_id uuid not null references author(author_id)
);

create table book_category (
    book_id uuid not null references book(book_id),
    category_id uuid not null references category(category_id)
);

insert into author(author_id, name) values ('f9ad2e1f-7511-40d4-a560-a0a6de712671', 'Rop Pike');

INSERT INTO book(book_id, name, author_id, updated_at) 
		VALUES ('5e2803bd-d7b1-4016-bb98-2905d98906b2', 'kema', 'f9ad2e1f-7511-40d4-a560-a0a6de712671', current_timestamp) RETURNING book_id

insert into category (category_id, name, parent_uuid) values
('9272dc0a-f27e-4fb9-8448-100851c1b374', 'kitob', null)
;

insert into category (category_id, name, parent_uuid) values
('39af40c6-7e8a-4495-b16e-a1c58da4cba3', 'badiiy', '9272dc0a-f27e-4fb9-8448-100851c1b374'),
('3bee22f5-96d5-4d6b-bd16-89454cb0cd18', 'melodrama', '9272dc0a-f27e-4fb9-8448-100851c1b374'),
('5152113a-03bd-4259-8fda-c9dbe677bffd', 'fantastik', '9272dc0a-f27e-4fb9-8448-100851c1b374')
;

insert into book_category(book_id, category_id) values
('5e2803bd-d7b1-4016-bb98-2905d98906b2', '39af40c6-7e8a-4495-b16e-a1c58da4cba3'),
('5e2803bd-d7b1-4016-bb98-2905d98906b2', '3bee22f5-96d5-4d6b-bd16-89454cb0cd18')
;

select
    array_agg(c.name)
from
    book_category as bc
join
    book as b using(book_id)
join
    category as c using(category_id)
group by b.name
;

SELECT 
    book_id, 
    name, 
    author_id, 
    created_at, 
    updated_at,
    (
            select
                array_agg(c.name)::varchar[]
            from
                book_category as bc
            join
                book as b using(book_id)
            join
                category as c using(category_id)
            where b.book_id = '5e2803bd-d7b1-4016-bb98-2905d98906b2'
            group by b.name
	)
	FROM book WHERE book_id='5e2803bd-d7b1-4016-bb98-2905d98906b2' AND deleted_at IS NULL
;


select
    c.name,
    c.parent_uuid,
    bc.created_at,
    bc.updated_at
from
    book_category as bc
join
    book as b using(book_id)
join
    category as c using(category_id)
where b.book_id = '5e2803bd-d7b1-4016-bb98-2905d98906b2'
;

select
    bc.book_id
from
    book_category as bc
group by bc.book_id;

UPDATE book_category SET deleted_at=current_timestamp 
WHERE book_id='e05f29bc-99a4-40d7-a7f0-3c30666334f1' AND category_id='5152113a-03bd-4259-8fda-c9dbe677bffd' AND deleted_at IS NULL

