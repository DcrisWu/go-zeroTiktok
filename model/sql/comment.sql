create table comment
(
    id         bigint unsigned not null,
    user_id    bigint unsigned not null,
    video_id   bigint unsigned not null,
    content    varchar(255)    not null,
    created_at datetime        not null default current_timestamp,
    updated_at datetime        not null default current_timestamp on update current_timestamp,
    deleted_at datetime                 default null,
    primary key (id),
    key video_id (video_id),
    unique key user_id_video_id (user_id, video_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;