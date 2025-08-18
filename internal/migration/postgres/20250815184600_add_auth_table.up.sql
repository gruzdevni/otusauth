CREATE TABLE logged_in(
    user_guid           UUID PRIMARY KEY        NOT NULL,
    expiry              TIMESTAMPTZ             
);

COMMENT ON COLUMN logged_in.user_guid       IS 'Гуид пользователя';
COMMENT ON COLUMN logged_in.expiry          IS 'Время окончания авторизации';


ALTER TABLE users ADD COLUMN pwd VARCHAR(255) NOT NULL DEFAULT ''::VARCHAR;
ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL DEFAULT ''::VARCHAR;

COMMENT ON COLUMN users.pwd IS 'Хэш пароля';
COMMENT ON COLUMN users.email IS 'Email';