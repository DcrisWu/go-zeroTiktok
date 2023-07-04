create table favorite
(
    id         bigint unsigned not null auto_increment,
    user_id    bigint unsigned not null,
    video_id   bigint unsigned not null,
    created_at datetime        not null default current_timestamp,
    updated_at datetime        not null default current_timestamp on update current_timestamp,
    primary key (id),
    unique key user_id_favorite_videos_id (user_id, video_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;