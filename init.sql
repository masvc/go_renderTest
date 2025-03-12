-- タスクのダミーデータ
INSERT INTO tasks (title, description, due_date, completed, created_at, updated_at) VALUES
('Renderへのデプロイ', 'Go WebアプリケーションをRenderにデプロイする', '2025-03-20', true, NOW(), NOW()),
('テストコードの作成', 'アプリケーションの単体テストとE2Eテストを実装する', '2025-03-25', false, NOW(), NOW()),
('ユーザー認証の実装', 'JWTを使用したユーザー認証システムを追加する', '2025-03-30', false, NOW(), NOW()),
('検索機能の追加', 'タスクのタイトルと説明で検索できる機能を実装する', '2025-04-05', false, NOW(), NOW()),
('UIの改善', 'モバイル対応とダークモードの実装', '2025-04-10', false, NOW(), NOW()); 