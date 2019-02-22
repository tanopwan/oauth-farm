-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table app_user (
	id				    bigserial	primary key   not null,
  username      text  not null,
	firstname			text,
	lastname			text,
	email         text  not null,
  status        int   not null default 0,
  lastLogin     timestamp with time zone	not null default now(),
  updated       timestamp with time zone	not null default now(),
	created       timestamp with time zone	not null default now(),

  unique (username),
  unique (email)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table app_user;