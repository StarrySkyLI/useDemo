create database demo;
use demo;

CREATE TABLE `game` (
                        `id` int unsigned NOT NULL AUTO_INCREMENT,
                        `name` varchar(20) NOT NULL,
                        `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
                        `data` varchar(20) NOT NULL,
                        `created_at` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP,
                        `updated_at` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP,
                        `deleted_at` timestamp NULL DEFAULT NULL,
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;