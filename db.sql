--Tables--

CREATE TABLE users(
    user_id     SERIAL,
    username    VARCHAR(100)    NOT NULL UNIQUE,
    password    VARCHAR(100)    NOT NULL UNIQUE,
    email       VARCHAR(100)    NOT NULL,
    role        VARCHAR(100)    NOT NULL,
    sowalga_tmt   INT           DEFAULT 0  NOT NULL,
    sowalga_usd   INT           DEFAULT 0  NOT NULL,
    is_it_deleted  VARCHAR(100) DEFAULT 'False'
);

CREATE TABLE workers(
    worker_id           SERIAL,
    fullname            VARCHAR(100)             NOT NULL,
    degisli_dukany      VARCHAR(100)             NOT NULL,
    wezipesi            VARCHAR(100)             NOT NULL,
    salary              INT                      NOT NULL,
    phone               VARCHAR(100)             NOT NULL,
    home_addres         VARCHAR(100)             NOT NULL,
    home_phone          VARCHAR(100)             NOT NULL,
    email               VARCHAR(100)             NOT NULL,   
    file_name           VARCHAR(100),
    note                VARCHAR(250),
    is_it_deleted       VARCHAR(100)             DEFAULT 'False'
);

CREATE TABLE stores(
    store_id                SERIAL,
    name                    VARCHAR(100)             NOT NULL,
    parent_store_id         INT    DEFAULT 1,       
    jemi_hasap_tmt          INT    DEFAULT 0         NOT NULL,
    jemi_hasap_usd          INT    DEFAULT 0         NOT NULL,
    shahsy_hasap_tmt        INT    DEFAULT 0         NOT NULL,
    shahsy_hasap_usd        INT    DEFAULT 0         NOT NULL,
    is_it_deleted           VARCHAR(100)             DEFAULT 'False'
);

CREATE TABLE categories(
    categorie_id       SERIAL,
    name               VARCHAR(100)     NOT NULL,
    parent_categorie   VARCHAR(100) DEFAULT '',
    is_it_deleted      VARCHAR(100) DEFAULT 'False'
);

CREATE TABLE customers(
    customer_id     SERIAL,
    name            VARCHAR(100)    NOT NULL,
    girdeyjisi_tmt      int         DEFAULT 0,
    girdeyjisi_usd      int         DEFAULT 0,
    is_it_deleted   VARCHAR(100)    DEFAULT 'False',
    note            VARCHAR(100) 
);

CREATE TABLE money_transfers(
    id                           SERIAL,
    store_id                     INT                NOT NULL,
    type_of_transfer             VARCHAR(100)       NOT NULL,        /* girdi, cykdy */
    type_of_account              VARCHAR(100)       NOT NULL,        /*nagt, bank*/
    currency                     VARCHAR(100)       NOT NULL,        
    categorie                    VARCHAR(100)       NOT NULL,
    date                         date,          
    keyword                      VARCHAR(255)       DEFAULT NULL,
    customer                     VARCHAR(100)       DEFAULT NULL,
    project                      VARCHAR(100)       DEFAULT NULL,
    type_of_income_payment       VARCHAR(100)       DEFAULT NULL,      /*doly, doly dal*/
    total_payment_amount         INT                DEFAULT NULL,      /*mocberi*/
    jemi_gecirilmelisi           INT                DEFAULT NULL,
    gecirileni                   INT                DEFAULT NULL,
    galany                       INT                DEFAULT NULL,
    money_gone_to                VARCHAR(100)       DEFAULT '',
    user_id                      INT                NOT NULL,
    note                         VARCHAR(255)       DEFAULT NULL,
    contract                     VARCHAR(255)       DEFAULT NULL,
    create_ts TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transfers_between_stores(
    id                     SERIAL,
    user_id                INT                    NOT NULL,
    from_store_name        VARCHAR(100)           NOT NULL,
    to_store_name          VARCHAR(100)           NOT NULL,
    type_of_account        VARCHAR(100)           NOT NULL,
    currency               VARCHAR(100)           NOT NULL,
    total_payment_amount   INT                    NOT NULL,
    date                   date                   NOT NULL,
    note                   VARCHAR(255),
    create_ts TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE last_modifications(
    id                       SERIAL,
    user_id                  INT                NOT NULL,
    action                   VARCHAR(255)       NOT NULL,
    message                  VARCHAR(255)       NOT NULL,
    is_it_seen               INT                DEFAULT 0,
    
    create_ts TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE income_outcome(
    id                 SERIAL,
    total_income_tmt   INT        DEFAULT 0,
    total_income_usd   INT        DEFAULT 0,
    total_outcome_tmt  INT        DEFAULT 0,
    total_outcome_usd  INT        DEFAULT 0

)













