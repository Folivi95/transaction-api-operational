/* When changing a SQL migration script, run
go-bindata -pkg migrations -ignore bindata -nometadata -prefix internal/adapters/postgres/migrations/
   -o ./internal/adapters/postgres/migrations/bindata.go ./internal/adapters/postgres/migrations
to update bindata.go
*/

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS transactions_api
(
    id          SERIAL      NOT NULL PRIMARY KEY,
    transaction jsonb       NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON transactions_api
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS acnt_contract
(
    id                 INT NOT NULL PRIMARY KEY,
    amnd_date          TIMESTAMP WITHOUT TIME ZONE,
    amnd_state         VARCHAR,
    amnd_officer       INT,
    amnd_prev          INT,
    pcat               VARCHAR,
    con_cat            VARCHAR,
    terminal_category  VARCHAR,
    ccat               VARCHAR,
    f_i                INT,
    branch             VARCHAR,
    service_group      VARCHAR,
    contract_number    VARCHAR,
    contract_name      VARCHAR,
    comment_text       VARCHAR,
    base_relation      VARCHAR,
    relation_tag       VARCHAR,
    acnt_contract__id  INT,
    contr_type         INT,
    contr_subtype__id  INT,
    serv_pack__id      INT,
    old_pack           INT,
    channel            VARCHAR,
    acc_scheme__id     INT,
    old_scheme         INT,
    product            VARCHAR,
    product_prev       VARCHAR,
    parent_product     VARCHAR,
    main_product       INT,
    client__id         INT,
    client_type        INT,
    acnt_contract__oid INT,
    liab_category      VARCHAR,
    liab_contract      INT,
    liab_contract_prev INT,
    billing_contract   INT,
    behavior_group     INT,
    behavior_type      INT,
    behavior_type_prev INT,
    check_available    VARCHAR,
    check_usage        VARCHAR,
    curr               VARCHAR,
    old_curr           VARCHAR,
    auth_limit_amount  DOUBLE PRECISION,
    base_auth_limit    DOUBLE PRECISION,
    liab_balance       DOUBLE PRECISION,
    liab_blocked       DOUBLE PRECISION,
    own_balance        DOUBLE PRECISION,
    own_blocked        DOUBLE PRECISION,
    sub_balance        DOUBLE PRECISION,
    sub_blocked        DOUBLE PRECISION,
    total_blocked      DOUBLE PRECISION,
    total_balance      DOUBLE PRECISION,
    shared_balance     DOUBLE PRECISION,
    shared_blocked     DOUBLE PRECISION,
    amount_available   DOUBLE PRECISION,
    date_open          TIMESTAMP WITHOUT TIME ZONE,
    date_expire        TIMESTAMP WITHOUT TIME ZONE,
    last_billing_date  TIMESTAMP WITHOUT TIME ZONE,
    next_billing_date  TIMESTAMP WITHOUT TIME ZONE,
    last_scan          TIMESTAMP WITHOUT TIME ZONE,
    card_expire        VARCHAR,
    production_status  VARCHAR,
    rbs_member_id      VARCHAR,
    rbs_number         VARCHAR,
    report_type        VARCHAR,
    max_pin_attempts   INT,
    pin_attempts       INT,
    chip_scheme        INT,
    risk_scheme        INT,
    risk_factor        DOUBLE PRECISION,
    risk_factor_prev   DOUBLE PRECISION,
    contr_status       INT,
    merchant_id        VARCHAR,
    tr_title           INT,
    tr_company         VARCHAR,
    tr_country         VARCHAR,
    tr_first_nam       VARCHAR,
    tr_last_nam        VARCHAR,
    tr_sic             VARCHAR,
    add_info_01        VARCHAR,
    add_info_02        VARCHAR,
    add_info_03        VARCHAR,
    add_info_04        VARCHAR,
    contract_level     VARCHAR,
    ext_data           VARCHAR,
    report_address     INT,
    share_balance      VARCHAR,
    is_multycurrency   VARCHAR,
    enables_item       VARCHAR,
    cycle_length       INT,
    interval_type      VARCHAR,
    status_category    VARCHAR,
    limit_is_active    VARCHAR,
    settlement_type    VARCHAR,
    is_ready           VARCHAR,
    auth_seq_n         INT,
    routing_idt        VARCHAR,
    apply_dt           TIMESTAMP WITHOUT TIME ZONE,
    local_version      INT,
    remote_version     INT
);

CREATE TABLE IF NOT EXISTS client_address
(
    id                 INT NOT NULL PRIMARY KEY,
    acnt_contract__oid INT,
    add_info           VARCHAR,
    address_line_1     VARCHAR,
    address_line_2     VARCHAR,
    address_line_3     VARCHAR,
    address_line_4     VARCHAR,
    address_name       VARCHAR,
    address_type       INT,
    address_zip        VARCHAR,
    amnd_date          TIMESTAMP WITHOUT TIME ZONE,
    amnd_officer       INT,
    amnd_prev          INT,
    amnd_state         VARCHAR,
    birth_nam          VARCHAR,
    city               VARCHAR,
    client__oid        INT,
    copy_to_address    INT,
    country            VARCHAR,
    date_from          TIMESTAMP WITHOUT TIME ZONE,
    date_to            TIMESTAMP WITHOUT TIME ZONE,
    delivery_type      INT,
    e_mail             VARCHAR,
    fax                VARCHAR,
    fax_h              VARCHAR,
    first_nam          VARCHAR,
    is_active          VARCHAR,
    is_ready           VARCHAR,
    language           INT,
    last_nam           VARCHAR,
    location           VARCHAR,
    municipality_code  VARCHAR,
    parent_address     INT,
    phone              VARCHAR,
    phone_h            VARCHAR,
    phone_m            VARCHAR,
    salutation_suffix  VARCHAR,
    state              VARCHAR,
    title              INT,
    url                VARCHAR,
    zip_code           VARCHAR
);

CREATE TABLE IF NOT EXISTS trans_type
(
    id                INT NOT NULL PRIMARY KEY,
    amnd_date         TIMESTAMP WITHOUT TIME ZONE,
    amnd_state        VARCHAR,
    amnd_prev         INT,
    amnd_officer      INT,
    name              VARCHAR,
    service_class     VARCHAR,
    s_cat             VARCHAR,
    t_cat             VARCHAR,
    dr_cr             INT,
    is_impersonal     VARCHAR,
    is_authorized     VARCHAR,
    is_required       VARCHAR,
    enable_adjustment VARCHAR,
    enable_reversal   VARCHAR,
    enable_request    VARCHAR,
    prev_trans_type   INT,
    chain_type        VARCHAR,
    charge_event      VARCHAR,
    dispute_trn_class VARCHAR,
    terminal_category VARCHAR,
    production_type   VARCHAR,
    production_event  VARCHAR,
    trans_code        VARCHAR,
    reversal_code     VARCHAR,
    trans_type_idt    VARCHAR,
    priority          INT
);

CREATE TABLE IF NOT EXISTS trans_cond
(
    id                INT NOT NULL PRIMARY KEY,
    amnd_state        VARCHAR,
    amnd_date         TIMESTAMP WITHOUT TIME ZONE,
    amnd_officer      INT,
    amnd_prev         INT,
    name              VARCHAR,
    code              VARCHAR,
    term_cat          VARCHAR,
    category_code     VARCHAR,
    condition_details VARCHAR,
    default_condition INT,
    late_condition    INT,
    security_code     VARCHAR,
    addendum          VARCHAR
);
