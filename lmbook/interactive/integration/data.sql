create database if not exists lmbook;
create table if not exists lmbook.interactives
(
    id          bigint auto_increment
    primary key,
    biz_id      bigint       null,
    biz         varchar(128) null,
    read_cnt    bigint       null,
    collect_cnt bigint       null,
    like_cnt    bigint       null,
    ctime       bigint       null,
    utime       bigint       null,
    constraint biz_type_id
    unique (biz_id, biz)
    );

create table if not exists lmbook.user_collection_bizs
(
    id     bigint auto_increment
    primary key,
    cid    bigint       null,
    biz_id bigint       null,
    biz    varchar(128) null,
    uid    bigint       null,
    ctime  bigint       null,
    utime  bigint       null,
    constraint biz_type_id_uid
    unique (biz_id, biz, uid)
    );

create index idx_user_collection_bizs_cid
    on lmbook.user_collection_bizs (cid);

create table if not exists lmbook.user_like_bizs
(
    id     bigint auto_increment
    primary key,
    biz_id bigint           null,
    biz    varchar(128)     null,
    uid    bigint           null,
    status tinyint unsigned null,
    ctime  bigint           null,
    utime  bigint           null,
    constraint biz_type_id_uid
    unique (biz_id, biz, uid)
    );

INSERT INTO `interactives`(`biz_id`, `biz`, `read_cnt`, `collect_cnt`, `like_cnt`, `ctime`, `utime`)
VALUES(1,"test",3128,6337,5817,1754998672039,1754998672039),
(2,"test",9066,5748,4033,1754998672039,1754998672039),
(3,"test",8686,4218,5507,1754998672039,1754998672039),
(4,"test",3809,3171,7878,1754998672039,1754998672039),
(5,"test",9773,7907,415,1754998672039,1754998672039),
(6,"test",3230,9390,5507,1754998672039,1754998672039),
(7,"test",1292,2600,9731,1754998672039,1754998672039),
(8,"test",4522,5635,4854,1754998672039,1754998672039),
(9,"test",4168,4215,44,1754998672039,1754998672039),
(10,"test",6085,6142,5349,1754998672039,1754998672039)