create table relation
(
    id         bigint unsigned not null,
    user_id    bigint unsigned not null,
    to_user_id bigint unsigned not null,
    created_at datetime        not null default current_timestamp,
    updated_at datetime        not null default current_timestamp on update current_timestamp,
    deleted_at datetime                 default null,
    primary key (id)
)