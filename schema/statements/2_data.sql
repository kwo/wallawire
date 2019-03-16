-- +migrate Up
UPSERT INTO roles (id, name) VALUES ('05ed6375-0786-4e1a-bcb0-1533c837954d', 'admin');
UPSERT INTO roles (id, name) VALUES ('ab9f2901-5aea-43b6-8f2b-7bf97dd30808', 'user');

UPSERT INTO users (id, disabled, username, name, created, updated, password_hash)
VALUES ('30d3bb1a-affa-48d6-859c-3c1537f9edbb', FALSE, 'demouser', 'Demo User', EXTRACT('epoch', now()), EXTRACT('epoch', now()), '243261243130245154466456787147626e4a4577665531556b5a7030657655456930583733492e366373516b667371696e30346c4c6d48327a727347');

UPSERT INTO user_role (user_id, role_id) VALUES ('30d3bb1a-affa-48d6-859c-3c1537f9edbb', '05ed6375-0786-4e1a-bcb0-1533c837954d');
UPSERT INTO user_role (user_id, role_id) VALUES ('30d3bb1a-affa-48d6-859c-3c1537f9edbb', 'ab9f2901-5aea-43b6-8f2b-7bf97dd30808');

-- +migrate Down
DELETE FROM user_role WHERE role_id IN ('05ed6375-0786-4e1a-bcb0-1533c837954d', 'ab9f2901-5aea-43b6-8f2b-7bf97dd30808');
DELETE FROM roles WHERE id IN ('05ed6375-0786-4e1a-bcb0-1533c837954d', 'ab9f2901-5aea-43b6-8f2b-7bf97dd30808');
