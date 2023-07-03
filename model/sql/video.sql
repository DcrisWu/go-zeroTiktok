create table video
(
    id             bigint unsigned not null,
    author_id      bigint unsigned not null,
    play_url       varchar(255)    not null,
    cover_url      varchar(255)    not null,
    favorite_count bigint unsigned not null default 0,
    comment_count  bigint unsigned not null default 0,
    title          varchar(50)     not null,
    created_at      datetime        not null default current_timestamp,
    updated_at      datetime        not null default current_timestamp on update current_timestamp,
    deleted_at      datetime                 default null,
    primary key (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;