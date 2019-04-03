drop table if exists `job_info`;
create table `job_info`(
    `id` int(10) unsigned not null auto_increment,
    `name` varchar(100) not null,
    `group` varchar(100) not null,
	`cron`  varchar(100) not null,
	`urls`  varchar(100) not null,
	`params` varchar(100) not null,
	`is_active` int(10) not null,
	`create_time` datetime not null,
	`modify_time` datetime not null,
    primary key(`id`)
)engine=InnoDB default charset=utf8;

drop table if exists `job_snapshot`;
create table `job_snapshot`(
	`id` int(10) unsigned not null auto_increment,
	`job_id` int(10) not null,
	`name` varchar(100) not null,
    `group` varchar(100) not null,
	`status`  varchar(100) not null,
	`url`  varchar(100) not null,
	`params` varchar(100) not null,
	`create_time` datetime not null,
	`modify_time` datetime not null,
	`time_consume` varchar(100) not null,
	primary key(`id`)
)engine=InnoDB default charset=utf8;

drop table if exists `user_info`;
create table `user_info`(
	`id` int(10) unsigned not null auto_increment,
	`username` varchar(100) not null,
	`password` varchar(100) not null,
	primary key(`id`),
	key(`username`)
)engine=InnoDB default charset=utf8;