ALTER TABLE author ADD COLUMN created_at timestamp default current_timestamp;
ALTER TABLE author ADD COLUMN updated_at timestamp;
ALTER TABLE author ADD COLUMN deleted_at timestamp;

ALTER TABLE category ADD COLUMN created_at timestamp default current_timestamp;
ALTER TABLE category ADD COLUMN updated_at timestamp;
ALTER TABLE category ADD COLUMN deleted_at timestamp;

ALTER TABLE book ADD COLUMN created_at timestamp default current_timestamp;
ALTER TABLE book ADD COLUMN updated_at timestamp;
ALTER TABLE book ADD COLUMN deleted_at timestamp;

ALTER TABLE book_category ADD COLUMN created_at timestamp default current_timestamp;
ALTER TABLE book_category ADD COLUMN updated_at timestamp;
ALTER TABLE book_category ADD COLUMN deleted_at timestamp;