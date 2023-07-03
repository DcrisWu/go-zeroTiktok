create table user
(
    id              bigint unsigned not null,
    user_name       varchar(40)     not null unique,
    password        varchar(256)    not null,
    following_count int unsigned             default 0,
    follower_count  int unsigned             default 0,
    created_at      datetime        not null default current_timestamp,
    updated_at      datetime        not null default current_timestamp on update current_timestamp,
    deleted_at      datetime                 default null,
    primary key (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;