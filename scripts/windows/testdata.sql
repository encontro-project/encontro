-- USERS
INSERT INTO users (id, username, email, password_hash)
VALUES (100, 'test_user', 'test_user@example.com', 'test_hash')
ON CONFLICT (id) DO NOTHING;

-- SERVERS
INSERT INTO servers (id, name, owner_id)
VALUES (1, 'Test Server', 100)
ON CONFLICT (id) DO NOTHING;

-- SERVER_USERS
INSERT INTO server_users (server_id, user_id, role)
VALUES (1, 100, 'owner')
ON CONFLICT (server_id, user_id) DO NOTHING;

-- CHATS
INSERT INTO chats (id, server_id, name, type)
VALUES
  (10, 1, 'General Chat', 'text'),
  (20, 1, 'General Voice', 'voice')
ON CONFLICT (id) DO NOTHING;

-- ROOMS
INSERT INTO rooms (id, name, type)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'General Chat Room', 'text'),
  ('22222222-2222-2222-2222-222222222222', 'General Voice Room', 'voice')
ON CONFLICT (id) DO NOTHING;

-- MESSAGES
INSERT INTO messages (id, room_id, content, sender_id)
VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'Hello, world!', 100),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 'Second message', 100)
ON CONFLICT (id) DO NOTHING;