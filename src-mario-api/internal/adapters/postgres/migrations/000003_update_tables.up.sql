DROP TABLE IF EXISTS acnt_contract;
DROP TABLE IF EXISTS client_address;
DROP TABLE IF EXISTS trans_type;
DROP TABLE IF EXISTS trans_cond;

CREATE TABLE IF NOT EXISTS acnt_contract
(
    id                 BIGINT NOT NULL PRIMARY KEY,
    amnd_date          TIMESTAMP WITHOUT TIME ZONE,
    amnd_state         VARCHAR,
    amnd_officer       BIGINT,
    amnd_prev          BIGINT,
    pcat               VARCHAR,
    con_cat            VARCHAR,
    terminal_category  VARCHAR,
    ccat               VARCHAR,
    f_i                BIGINT,
    branch             VARCHAR,
    service_group      VARCHAR,
    contract_number    VARCHAR,
    contract_name      VARCHAR,
    comment_text       VARCHAR,
    base_relation      VARCHAR,
    relation_tag       VARCHAR,
    acnt_contract__id  BIGINT,
    contr_type         BIGINT,
    contr_subtype__id  BIGINT,
    serv_pack__id      BIGINT,
    old_pack           BIGINT,
    channel            VARCHAR,
    acc_scheme__id     BIGINT,
    old_scheme         BIGINT,
    product            VARCHAR,
    product_prev       VARCHAR,
    parent_product     VARCHAR,
    main_product       BIGINT,
    client__id         BIGINT,
    client_type        BIGINT,
    acnt_contract__oid BIGINT,
    liab_category      VARCHAR,
    liab_contract      BIGINT,
    liab_contract_prev BIGINT,
    billing_contract   BIGINT,
    behavior_group     BIGINT,
    behavior_type      BIGINT,
    behavior_type_prev BIGINT,
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
    max_pin_attempts   BIGINT,
    pin_attempts       BIGINT,
    chip_scheme        BIGINT,
    risk_scheme        BIGINT,
    risk_factor        DOUBLE PRECISION,
    risk_factor_prev   DOUBLE PRECISION,
    contr_status       BIGINT,
    merchant_id        VARCHAR,
    tr_title           BIGINT,
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
    report_address     BIGINT,
    share_balance      VARCHAR,
    is_multycurrency   VARCHAR,
    enables_item       VARCHAR,
    cycle_length       BIGINT,
    interval_type      VARCHAR,
    status_category    VARCHAR,
    limit_is_active    VARCHAR,
    settlement_type    VARCHAR,
    is_ready           VARCHAR,
    auth_seq_n         BIGINT,
    routing_idt        VARCHAR,
    apply_dt           TIMESTAMP WITHOUT TIME ZONE,
    local_version      BIGINT,
    remote_version     BIGINT
);

CREATE TABLE IF NOT EXISTS client_address
(
    id                 BIGINT NOT NULL PRIMARY KEY,
    acnt_contract__oid BIGINT,
    add_info           VARCHAR,
    address_line_1     VARCHAR,
    address_line_2     VARCHAR,
    address_line_3     VARCHAR,
    address_line_4     VARCHAR,
    address_name       VARCHAR,
    address_type       BIGINT,
    address_zip        VARCHAR,
    amnd_date          TIMESTAMP WITHOUT TIME ZONE,
    amnd_officer       BIGINT,
    amnd_prev          BIGINT,
    amnd_state         VARCHAR,
    birth_nam          VARCHAR,
    city               VARCHAR,
    client__oid        BIGINT,
    copy_to_address    BIGINT,
    country            VARCHAR,
    date_from          TIMESTAMP WITHOUT TIME ZONE,
    date_to            TIMESTAMP WITHOUT TIME ZONE,
    delivery_type      BIGINT,
    e_mail             VARCHAR,
    fax                VARCHAR,
    fax_h              VARCHAR,
    first_nam          VARCHAR,
    is_active          VARCHAR,
    is_ready           VARCHAR,
    language           BIGINT,
    last_nam           VARCHAR,
    location           VARCHAR,
    municipality_code  VARCHAR,
    parent_address     BIGINT,
    phone              VARCHAR,
    phone_h            VARCHAR,
    phone_m            VARCHAR,
    salutation_suffix  VARCHAR,
    state              VARCHAR,
    title              BIGINT,
    url                VARCHAR,
    zip_code           VARCHAR
);

CREATE TABLE IF NOT EXISTS trans_type
(
    id                BIGINT NOT NULL PRIMARY KEY,
    amnd_date         TIMESTAMP WITHOUT TIME ZONE,
    amnd_state        VARCHAR,
    amnd_prev         BIGINT,
    amnd_officer      BIGINT,
    name              VARCHAR,
    service_class     VARCHAR,
    s_cat             VARCHAR,
    t_cat             VARCHAR,
    dr_cr             BIGINT,
    is_impersonal     VARCHAR,
    is_authorized     VARCHAR,
    is_required       VARCHAR,
    enable_adjustment VARCHAR,
    enable_reversal   VARCHAR,
    enable_request    VARCHAR,
    prev_trans_type   BIGINT,
    chain_type        VARCHAR,
    charge_event      VARCHAR,
    dispute_trn_class VARCHAR,
    terminal_category VARCHAR,
    production_type   VARCHAR,
    production_event  VARCHAR,
    trans_code        VARCHAR,
    reversal_code     VARCHAR,
    trans_type_idt    VARCHAR,
    priority          BIGINT
);

CREATE TABLE IF NOT EXISTS trans_cond
(
    id                BIGINT NOT NULL PRIMARY KEY,
    amnd_state        VARCHAR,
    amnd_date         TIMESTAMP WITHOUT TIME ZONE,
    amnd_officer      BIGINT,
    amnd_prev         BIGINT,
    name              VARCHAR,
    code              VARCHAR,
    term_cat          VARCHAR,
    category_code     VARCHAR,
    condition_details VARCHAR,
    default_condition BIGINT,
    late_condition    BIGINT,
    security_code     VARCHAR,
    addendum          VARCHAR
);

CREATE TABLE IF NOT EXISTS aux_transactions_api
(
    id            BIGINT NOT NULL PRIMARY KEY,
    salt_token_id VARCHAR,
    expiry_date   VARCHAR
);

CREATE TABLE IF NOT EXISTS bin_table
(
    id         BIGINT NOT NULL PRIMARY KEY,
    country    VARCHAR,
    card_brand VARCHAR
);
