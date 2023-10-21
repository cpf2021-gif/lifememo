create database lifememo_user;
use lifememo_user;

CREATE TABLE `user` (
                        `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                        `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
                        `email` varchar(320) NOT NULL DEFAULT '' COMMENT '邮箱',
                        `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                        PRIMARY KEY (`id`),
                        KEY `ix_update_time` (`update_time`),
                        UNIQUE KEY `uk_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户表';