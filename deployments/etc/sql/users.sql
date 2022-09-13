CREATE TABLE `users` (
    `id_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `secret` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    PRIMARY KEY (`id_user`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

INSERT INTO `users`(`id_user`, `secret`) VALUES ('admin', '$2a$12$451sd5GjT4uEAvPTtkLIdOl1KdDLNSvKRkqjLFSo9i5BDIaggL10q');

CREATE TABLE `user_roles` (
    `id_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `id_role` smallint DEFAULT NULL,
    UNIQUE KEY `id_user` (`id_user`,`id_role`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

INSERT INTO `user_roles`(`id_user`, `id_role`) VALUES ('admin', 2);

CREATE TABLE `jwt_access_tokens` (
    `id_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `access_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    `expired_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;