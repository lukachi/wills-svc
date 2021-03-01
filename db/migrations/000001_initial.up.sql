CREATE TABLE accounts
(
    id          bigserial primary key,
    account_id  text        not null,
    role        smallserial not null,
    first_name  text        not null,
    middle_name text        not null,
    last_name   text        not null,
    password    text        not null,
    salt        text        not null,
    pk          text        not null
);

CREATE TABLE wills
(
    id           bigserial primary key,
    creator_id   bigserial references accounts (id) on delete cascade not null,
    recipient_id bigserial references accounts (id) on delete cascade not null,
    file_key     text                                                 not null
);