create database lifememo_message;
use lifememo_message;

CREATE TABLE `message` (
                          `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                          `send_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '发送消息用户ID',
                          `receive_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '接收消息用户ID',
                          `content` text COLLATE utf8_unicode_ci NOT NULL COMMENT '内容',
                          `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:审核中 1:正常 2:删除',
                          `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                          PRIMARY KEY (`id`),
                          KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='消息表';