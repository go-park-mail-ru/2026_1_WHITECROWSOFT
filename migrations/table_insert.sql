INSERT INTO account (id, username, password_hash, token_version, created_at, updated_at) VALUES
(1, 'john_doe', 'password123', 1, '2024-01-01 10:00:00+00', '2024-01-01 10:00:00+00'),
(2, 'jane_smith', 'securepass456', 1, '2024-01-02 11:30:00+00', '2024-01-02 11:30:00+00'),
(3, 'alex_wilson', 'strongpass789', 1, '2024-01-03 09:15:00+00', '2024-01-03 09:15:00+00'),
(4, 'emily_brown', 'pass12345678', 1, '2024-01-04 14:20:00+00', '2024-01-04 14:20:00+00'),
(5, 'mike_jones', 'mikepass1234', 2, '2024-01-05 16:45:00+00', '2024-01-05 16:45:00+00');


INSERT INTO block_type (id, name) VALUES
(1, 'text'),
(2, 'header'),
(3, 'code'),
(4, 'quote'),
(5, 'list'),
(6, 'image'),
(7, 'video'),
(8, 'table'),
(9, 'checklist');


INSERT INTO attachment (id, user_id, file_name, file_size, mime_type, storage_url, created_at) VALUES
(1, 1, 'profile-pic.jpg', 1024576, 'image/jpeg', 'https://storage.example.com/users/1/profile-pic.jpg', '2024-01-06 10:00:00+00'),
(2, 2, 'pasta-carbonara.jpg', 2048576, 'image/jpeg', 'https://storage.example.com/users/2/pasta-carbonara.jpg', '2024-01-07 11:30:00+00'),
(3, 3, 'go-cheatsheet.pdf', 512000, 'application/pdf', 'https://storage.example.com/users/3/go-cheatsheet.pdf', '2024-01-08 09:15:00+00'),
(4, 4, 'japan-travel-guide.pdf', 1024000, 'application/pdf', 'https://storage.example.com/users/4/japan-travel-guide.pdf', '2024-01-09 14:20:00+00'),
(5, 5, 'books-list.txt', 1024, 'text/plain', 'https://storage.example.com/users/5/books-list.txt', '2024-01-10 16:45:00+00');


INSERT INTO note (id, user_id, title, parent_id, icon_emoji, cover_picture_url, is_deleted, created_at, updated_at) VALUES
(1, 1, 'Личный дневник', NULL, '📔', NULL, false, '2024-01-10 09:00:00+00', '2024-01-10 09:00:00+00'),
(2, 1, 'Рабочие проекты', NULL, '💼', NULL, false, '2024-01-11 10:30:00+00', '2024-01-11 10:30:00+00'),
(3, 1, 'Заметки с митинга', 2, NULL, NULL, false, '2024-01-12 14:15:00+00', '2024-01-12 14:15:00+00'),
(4, 2, 'Кулинарная книга', NULL, '🍳', 'https://storage.example.com/covers/cooking.jpg', false, '2024-01-13 11:00:00+00', '2024-01-13 11:00:00+00'),
(5, 2, 'Любимые рецепты', 4, '🍝', NULL, false, '2024-01-14 12:20:00+00', '2024-01-14 12:20:00+00'),
(6, 3, 'Изучение Go', NULL, '🐹', NULL, false, '2024-01-15 15:30:00+00', '2024-01-15 15:30:00+00'),
(7, 4, 'Планы путешествий', NULL, '✈️', 'https://storage.example.com/covers/travel.jpg', false, '2024-01-16 08:45:00+00', '2024-01-16 08:45:00+00'),
(8, 5, 'Список книг для чтения', NULL, NULL, NULL, false, '2024-01-17 19:00:00+00', '2024-01-17 19:00:00+00');


INSERT INTO block (id, note_id, position, block_type_id, content, attachment_id, created_at, updated_at) VALUES
(1, 1, 0, 2, 'Мой первый день в новом году', NULL, '2024-01-10 09:01:00+00', '2024-01-10 09:01:00+00'),
(2, 1, 1, 1, 'Сегодня был отличный день! Начал заниматься спортом и прочитал 50 страниц книги.', NULL, '2024-01-10 09:02:00+00', '2024-01-10 09:02:00+00'),
(3, 1, 2, 9, 'Пробежка утром', NULL, '2024-01-10 09:03:00+00', '2024-01-10 09:03:00+00'),
(4, 2, 0, 2, 'Текущие проекты 2024', NULL, '2024-01-11 10:31:00+00', '2024-01-11 10:31:00+00'),
(5, 2, 1, 5, 'Разработка нового API', NULL, '2024-01-11 10:32:00+00', '2024-01-11 10:32:00+00'),
(6, 2, 2, 5, 'Рефакторинг базы данных', NULL, '2024-01-11 10:33:00+00', '2024-01-11 10:33:00+00'),
(7, 3, 0, 2, 'Еженедельный митинг 12.01', NULL, '2024-01-12 14:16:00+00', '2024-01-12 14:16:00+00'),
(8, 3, 1, 5, 'Обсуждение архитектуры микросервисов', NULL, '2024-01-12 14:17:00+00', '2024-01-12 14:17:00+00'),
(9, 4, 0, 2, 'Основные блюда', NULL, '2024-01-13 11:01:00+00', '2024-01-13 11:01:00+00'),
(10, 4, 1, 4, 'Хорошая еда - это основа счастливой жизни', NULL, '2024-01-13 11:02:00+00', '2024-01-13 11:02:00+00'),
(11, 5, 0, 2, 'Паста карбонара', NULL, '2024-01-14 12:21:00+00', '2024-01-14 12:21:00+00'),
(12, 5, 1, 5, 'Спагетти - 200г', NULL, '2024-01-14 12:22:00+00', '2024-01-14 12:22:00+00'),
(13, 5, 2, 5, 'Яйца - 3 шт', NULL, '2024-01-14 12:23:00+00', '2024-01-14 12:23:00+00'),
(14, 5, 3, 6, 'Фото готового блюда', 2, '2024-01-14 12:24:00+00', '2024-01-14 12:24:00+00'),
(15, 6, 0, 2, 'Основы Go', NULL, '2024-01-15 15:31:00+00', '2024-01-15 15:31:00+00'),
(16, 6, 1, 3, 'package main\n\nimport "fmt"\n\nfunc main() {\n    fmt.Println("Hello, World!")\n}', NULL, '2024-01-15 15:32:00+00', '2024-01-15 15:32:00+00'),
(17, 6, 2, 7, 'Видеоурок по Go', 3, '2024-01-15 15:33:00+00', '2024-01-15 15:33:00+00'),
(18, 7, 0, 2, 'Поездка в Японию 2024', NULL, '2024-01-16 08:46:00+00', '2024-01-16 08:46:00+00'),
(19, 7, 1, 5, 'Токио - 3 дня', NULL, '2024-01-16 08:47:00+00', '2024-01-16 08:47:00+00'),
(20, 7, 2, 5, 'Киото - 2 дня', NULL, '2024-01-16 08:48:00+00', '2024-01-16 08:48:00+00'),
(21, 7, 3, 7, 'Видео о Японии', 4, '2024-01-16 08:49:00+00', '2024-01-16 08:49:00+00'),
(22, 8, 0, 2, 'Техническая литература', NULL, '2024-01-17 19:01:00+00', '2024-01-17 19:01:00+00'),
(23, 8, 1, 1, 'Clean Code - Роберт Мартин', NULL, '2024-01-17 19:02:00+00', '2024-01-17 19:02:00+00'),
(24, 8, 2, 1, 'The Pragmatic Programmer', NULL, '2024-01-17 19:03:00+00', '2024-01-17 19:03:00+00'),
(25, 8, 3, 8, 'Список книг', 5, '2024-01-17 19:04:00+00', '2024-01-17 19:04:00+00');


INSERT INTO block_state (id, block_id, formatting, created_at) VALUES
(1, 1, '{"bold": true, "fontSize": 28, "color": "#2c3e50"}', '2024-01-10 09:01:00+00'),
(2, 2, '{"fontFamily": "Arial", "lineHeight": 1.6, "color": "#34495e"}', '2024-01-10 09:02:00+00'),
(3, 3, '{"checked": true, "bulletStyle": "checkbox"}', '2024-01-10 09:03:00+00'),
(4, 4, '{"bold": true, "fontSize": 24, "alignment": "center"}', '2024-01-11 10:31:00+00'),
(5, 7, '{"bold": true, "backgroundColor": "#f1c40f"}', '2024-01-12 14:16:00+00'),
(6, 10, '{"italic": true, "author": "Джулия Чайлд", "borderLeft": "3px solid #e67e22"}', '2024-01-13 11:02:00+00'),
(7, 11, '{"bold": true, "color": "#27ae60"}', '2024-01-14 12:21:00+00'),
(8, 14, '{"width": "100%", "borderRadius": "8px", "caption": "Готовая паста карбонара"}', '2024-01-14 12:24:00+00'),
(9, 16, '{"language": "go", "theme": "monokai", "lineNumbers": true, "highlightLines": [5]}', '2024-01-15 15:32:00+00'),
(10, 18, '{"bold": true, "fontSize": 26, "color": "#8e44ad"}', '2024-01-16 08:46:00+00'),
(11, 21, '{"width": "560", "height": "315", "autoplay": false}', '2024-01-16 08:49:00+00'),
(12, 22, '{"bold": true, "color": "#2980b9"}', '2024-01-17 19:01:00+00'),
(13, 23, '{"fontStyle": "italic", "color": "#16a085"}', '2024-01-17 19:02:00+00');


INSERT INTO favorite_note (id, user_id, note_id, added_at) VALUES
(1, 1, 4, '2024-01-15 10:00:00+00'),
(2, 1, 6, '2024-01-16 11:30:00+00'),
(3, 2, 3, '2024-01-17 09:15:00+00'),
(4, 3, 8, '2024-01-18 14:20:00+00'),
(5, 4, 2, '2024-01-19 16:45:00+00'),
(6, 5, 1, '2024-01-20 12:00:00+00');


INSERT INTO note_share (id, note_id, user_id, permission, shared_at, shared_by_user_id) VALUES
(1, 2, 2, 'R', '2024-01-15 10:00:00+00', 1),
(2, 3, 3, 'W', '2024-01-16 11:30:00+00', 1),
(3, 5, 1, 'R', '2024-01-17 09:15:00+00', 2),
(4, 6, 4, 'R', '2024-01-18 14:20:00+00', 3),
(5, 8, 2, 'W', '2024-01-19 16:45:00+00', 5);


SELECT setval('account_id_seq', (SELECT MAX(id) FROM account));
SELECT setval('note_id_seq', (SELECT MAX(id) FROM note));
SELECT setval('block_id_seq', (SELECT MAX(id) FROM block));
SELECT setval('block_state_id_seq', (SELECT MAX(id) FROM block_state));
SELECT setval('attachment_id_seq', (SELECT MAX(id) FROM attachment));
SELECT setval('note_share_id_seq', (SELECT MAX(id) FROM note_share));
SELECT setval('favorite_note_id_seq', (SELECT MAX(id) FROM favorite_note));
