CREATE TABLE users(
    guid                UUID PRIMARY KEY        NOT NULL,
    name                VARCHAR(255)            NOT NULL,
    occupation          TEXT                    NOT NULL,
    is_deleted          BOOLEAN                 NOT NULL DEFAULT false,
    created_at          TIMESTAMPTZ             NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ             NOT NULL DEFAULT now()
);

COMMENT ON COLUMN users.guid          IS 'GUID пользователя';
COMMENT ON COLUMN users.name          IS 'Имя пользователя';
COMMENT ON COLUMN users.occupation    IS 'Место работы';
COMMENT ON COLUMN users.is_deleted    IS 'Признак удален ли пользователь';
COMMENT ON COLUMN users.created_at    IS 'Дата создания';
COMMENT ON COLUMN users.updated_at    IS 'Дата обновления';

INSERT INTO users (
        guid,
        name,
        occupation,
        is_deleted,
        created_at,
        updated_at
    )
VALUES (
        '149497f4-aaf0-4881-86c7-498d191d3717',
        'Иванова Ариадна Евгеньевна',
        'МУП ДЭС',
        false,
        now(),
        now()
    ),
    (
        'f531286c-7d8f-4fd0-9900-8da398e371b5',
        'Степанов Эдуард',
        'МУП ДЭС',
        false,
        now(),
        now()
    ),
    (
        '8bff61fe-c8a1-45c7-895a-b0907c390279',
        'Сидоренко Валентин',
        'МУП ДЭС',
        false,
        now(),
        now()
    );