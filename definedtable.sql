CREATE TABLE account_login
(
    id            UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    account       VARCHAR(255) NOT NULL UNIQUE,
    account_type  VARCHAR(255) NOT NULL UNIQUE, -- 'business , internal'
    password_salt VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    created_by    UUID         NOT NULL,
    updated_by    UUID         NOT NULL
);

CREATE TABLE user_account
(
    id            UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    user_login_id UUID      NOT NULL REFERENCES account_login (id) UNIQUE,

    full_name     VARCHAR(255),
    avatar_url    TEXT,
    phone         VARCHAR(20),
    date_of_birth DATE,
    gender        VARCHAR(10),

    address       TEXT,
    is_active     BOOLEAN   NOT NULL DEFAULT TRUE,

    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by    UUID      NOT NULL,
    updated_by    UUID      NOT NULL
);

CREATE TABLE admin_account
(
    id               UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    account_login_id UUID      NOT NULL REFERENCES user_login (id) UNIQUE,

    full_name        VARCHAR(255),
    avatar_url       TEXT,
    phone            VARCHAR(20),
    date_of_birth    DATE,
    gender           VARCHAR(10),

    address          TEXT,
    is_active        BOOLEAN   NOT NULL DEFAULT TRUE,

    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by       UUID      NOT NULL,
    updated_by       UUID      NOT NULL
);

CREATE TABLE role
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name         VARCHAR(100) NOT NULL,
    role_type    INT          NOT NULL,
    account_type VARCHAR(100) NOT NULL, -- ('business', 'internal')
    description  TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- cronjob
CREATE TABLE account_zalo_oa
(
    id               UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    oa_id            VARCHAR(100) NOT NULL UNIQUE, -- ID từ Zalo
    oa_name          VARCHAR(255) NOT NULL,
    access_token     TEXT,                         -- sống
    refresh_token    TEXT,
    token_expires_at TIMESTAMPTZ,
    is_active        BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_by       UUID         NOT NULL,
    updated_by       UUID         NOT NULL
);

CREATE TABLE user_account_oa_mapping
(
    id              UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_account_id UUID        NOT NULL REFERENCES user_account (id),
    zalo_oa_id      UUID        NOT NULL REFERENCES account_zalo_oa (id),
    is_active       BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID        NOT NULL,
    updated_by      UUID        NOT NULL
        UNIQUE (user_account_id, zalo_oa_id)
);

CREATE TABLE user_internal_role
(
    id              UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    user_account_id UUID      NOT NULL REFERENCES user_account (id) UNIQUE,
    role_id         UUID      NOT NULL REFERENCES role (id),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by      BIGINT    NOT NULL,
    updated_by      BIGINT    NOT NULL created_by    UUID       NOT NULL,
    updated_by      UUID      NOT NULL
);

