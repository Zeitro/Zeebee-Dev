CREATE TABLE public.borrower_applicants
(
    id                 serial4      NOT NULL,
    created_at         timestamptz  NOT NULL DEFAULT now(),
    updated_at         timestamptz  NOT NULL DEFAULT now(),
    is_del             int2         NOT NULL DEFAULT 0,
    borrower_id        int4         NOT NULL DEFAULT 0,
    first_name         varchar(100) NOT NULL DEFAULT ''::character varying,
    mid_name           varchar(100) NOT NULL DEFAULT ''::character varying,
    last_name          varchar(100) NOT NULL DEFAULT ''::character varying,
    full_name          varchar(255) NOT NULL DEFAULT ''::character varying,
    ssn                varchar(20)  NOT NULL DEFAULT ''::character varying,
    dob                date NULL,
    marital_status     varchar(20)  NOT NULL DEFAULT ''::character varying,
    email              varchar(100) NOT NULL DEFAULT ''::character varying,
    est_credit_score   int2         NOT NULL DEFAULT 0,
    driver_license_num varchar(30)  NOT NULL DEFAULT ''::character varying,
    address            json         NOT NULL DEFAULT '{}'::json,
    phone              varchar(20)  NOT NULL DEFAULT ''::character varying,
    CONSTRAINT borrower_applicants_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_borrower_applicants_borrower_id ON public.borrower_applicants USING btree (borrower_id);


CREATE TABLE public.borrower_entities
(
    id              serial4      NOT NULL,
    created_at      timestamptz  NOT NULL DEFAULT now(),
    updated_at      timestamptz  NOT NULL DEFAULT now(),
    is_del          int2         NOT NULL DEFAULT 0,
    borrower_id     int4         NOT NULL DEFAULT 0,
    "type"          varchar(50)  NOT NULL DEFAULT ''::character varying,
    "name"          varchar(100) NOT NULL DEFAULT ''::character varying,
    tin_type        varchar(50)  NOT NULL DEFAULT ''::character varying,
    tin             varchar(20)  NOT NULL DEFAULT ''::character varying,
    mailing_address json         NOT NULL DEFAULT '{}'::json,
    contact_email   varchar(100) NOT NULL DEFAULT ''::character varying,
    contact_name    varchar(100) NOT NULL DEFAULT ''::character varying,
    contact_phone   varchar(20)  NOT NULL DEFAULT ''::character varying,
    CONSTRAINT borrower_entities_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_borrower_entities_borrower_id ON public.borrower_entities USING btree (borrower_id);


CREATE TABLE public.borrower_properties
(
    id               int4           NOT NULL DEFAULT nextval('properties_id_seq'::regclass),
    unit             int2           NOT NULL DEFAULT 1, -- number of unit
    created_at       timestamptz(6) NOT NULL DEFAULT now(),
    updated_at       timestamptz(6) NOT NULL DEFAULT now(),
    is_del           int2           NOT NULL DEFAULT 0,
    address          jsonb          NOT NULL DEFAULT '{}'::jsonb,
    borrower_id      int4           NOT NULL DEFAULT 0,
    total_amount     numeric(10, 2) NOT NULL DEFAULT 0, -- total loan amount
    property_type    varchar(50)    NOT NULL DEFAULT ''::character varying,
    living_area_sqft numeric(8, 2)  NOT NULL DEFAULT 0,
    bedrooms         int2           NOT NULL DEFAULT 0,
    bathrooms        int2           NOT NULL DEFAULT 0,
    loan_count       int4           NOT NULL DEFAULT 0,
    CONSTRAINT properties_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_property_borrower_id ON public.borrower_properties USING btree (borrower_id);

-- Column comments
COMMENT
ON COLUMN public.borrower_properties.unit IS 'number of unit';
COMMENT
ON COLUMN public.borrower_properties.total_amount IS 'total loan amount';


CREATE TABLE public.borrowers
(
    id         serial4      NOT NULL,
    user_id    int4         NOT NULL DEFAULT 0,
    first_name varchar(100) NOT NULL DEFAULT ''::character varying,
    last_name  varchar(100) NOT NULL DEFAULT ''::character varying,
    full_name  varchar(100) NOT NULL DEFAULT ''::character varying,
    created_at timestamp(6) NULL DEFAULT now(),
    updated_at timestamp(6) NULL DEFAULT now(),
    is_del     int2         NOT NULL DEFAULT 0,
    CONSTRAINT borrowers_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_borrower_user_id ON public.borrowers USING btree (user_id);


CREATE TABLE public.files
(
    created_at  timestamptz(6) NOT NULL DEFAULT now(),
    updated_at  timestamptz(6) NOT NULL DEFAULT now(),
    is_del      int2         NOT NULL DEFAULT 0,
    entity_id   int4         NOT NULL DEFAULT 0,
    file_url    varchar(255) NOT NULL DEFAULT ''::character varying,
    id          serial4      NOT NULL,
    entity_type varchar(50)  NOT NULL DEFAULT ''::character varying
);
CREATE INDEX idx_files_ref_id ON public.files USING btree (entity_id);


CREATE TABLE public.lender_companies
(
    id            int4         NOT NULL DEFAULT nextval('lenders_id_seq'::regclass),
    "name"        varchar(100) NOT NULL DEFAULT ''::character varying,
    status        varchar(100) NOT NULL DEFAULT 'active'::character varying,     -- active,inactive
    created_at    timestamptz(6) NOT NULL DEFAULT now(),
    updated_at    timestamptz(6) NOT NULL DEFAULT now(),
    is_del        int2         NOT NULL DEFAULT 0,
    address       jsonb        NOT NULL DEFAULT '{}'::jsonb,
    contact_email varchar(255) NOT NULL DEFAULT ''::character varying,
    owner_type    varchar(30)  NOT NULL DEFAULT 'ThirdParty'::character varying, -- Self,Third Party,Partner
    CONSTRAINT lenders_pkey PRIMARY KEY (id)
);

-- Column comments
COMMENT
ON COLUMN public.lender_companies.status IS 'active,inactive';
COMMENT
ON COLUMN public.lender_companies.owner_type IS 'Self,Third Party,Partner';


INSERT INTO public.lender_companies ("name", status, created_at, updated_at, is_del, address, contact_email, owner_type)
VALUES ('zeebee', 'Active', now(), now(), 0, '{}'::jsonb, 'zeebeelender@zeebeeloans.com', 'Self'),
       ('RBI Private Lending', 'Active', now(), now(), 0, '{}'::jsonb, 'rbiloan@zeebeeloans.com', 'ThirdParty'),
       ('Anchor Loans', 'Active', now(), now(), 0, '{}'::jsonb, 'anchorloan@zeebeeloans.com', 'ThirdParty');



CREATE TABLE public.lender_company_and_employees
(
    id                 serial4 NOT NULL,
    lender_company_id  int4    NOT NULL DEFAULT 0,
    lender_employee_id int4    NOT NULL DEFAULT 0,
    CONSTRAINT lender_company_and_employees_pkey PRIMARY KEY (id)
);
CREATE INDEX "idx_lender_company_relation_employee_id" ON public.lender_company_and_employees USING btree (lender_employee_id);
CREATE INDEX idx_lender_company_relation_company_id ON public.lender_company_and_employees USING btree (lender_company_id);


CREATE TABLE public.lender_employees
(
    id                        int4         NOT NULL DEFAULT nextval('borrowers_id_seq'::regclass),
    user_id                   int4         NOT NULL DEFAULT 0,
    first_name                varchar(100) NOT NULL DEFAULT ''::character varying,
    last_name                 varchar(100) NOT NULL DEFAULT ''::character varying,
    full_name                 varchar(100) NOT NULL DEFAULT ''::character varying,
    created_at                timestamp(6) NULL DEFAULT now(),
    updated_at                timestamp(6) NULL DEFAULT now(),
    is_del                    int2         NOT NULL DEFAULT 0,
    default_lender_company_id int4         NOT NULL DEFAULT 0,
    CONSTRAINT borrowers_copy1_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_default_lender_company_id ON public.lender_employees USING btree (default_lender_company_id);
CREATE INDEX idx_lender_user_id ON public.lender_employees USING btree (user_id);

CREATE TABLE public.loan_applicants
(
    id                    int4         NOT NULL DEFAULT nextval('applicants_id_seq'::regclass),
    full_name             varchar(255) NOT NULL DEFAULT ''::character varying,
    created_at            timestamptz(6) NOT NULL DEFAULT now(),
    updated_at            timestamptz(6) NOT NULL DEFAULT now(),
    is_del                int2         NOT NULL DEFAULT 0,
    borrower_id           int4         NOT NULL DEFAULT 0,
    loan_id               int4         NOT NULL DEFAULT 0,
    is_primary            bool         NOT NULL DEFAULT false,
    loan_property_id      int4         NOT NULL DEFAULT 0,
    "type"                varchar(10)  NOT NULL DEFAULT 'Individual'::character varying, -- Individual,Company
    invitation_status     varchar(10)  NOT NULL DEFAULT 'Pending'::character varying,    -- Accepted,Pending,Rejected
    ssn                   varchar(20)  NOT NULL DEFAULT ''::character varying,
    dob                   date NULL,
    marital_status        varchar(10)  NOT NULL DEFAULT ''::character varying,
    email                 varchar(100) NOT NULL DEFAULT ''::character varying,
    est_credit_score      int2         NOT NULL DEFAULT 0,
    driver_license_num    varchar(20)  NOT NULL DEFAULT ''::character varying,
    address               jsonb        NOT NULL DEFAULT '{}'::jsonb,
    loan_borrower_id      int4         NOT NULL DEFAULT 0,
    first_name            varchar(100) NOT NULL DEFAULT ''::character varying,
    mid_name              varchar(30)  NOT NULL DEFAULT ''::character varying,
    last_name             varchar(100) NOT NULL DEFAULT ''::character varying,
    phone                 varchar(20)  NOT NULL DEFAULT ''::character varying,
    borrower_applicant_id int4         NOT NULL DEFAULT 0,
    own_percent           int2         NOT NULL DEFAULT 0,
    CONSTRAINT applicants_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_applicant_borrower_id ON public.loan_applicants USING btree (borrower_id);
CREATE INDEX idx_applicant_loan_bo_id ON public.loan_applicants USING btree (loan_borrower_id);
CREATE INDEX idx_applicant_loan_id ON public.loan_applicants USING btree (loan_id);
CREATE INDEX idx_applicant_loan_property_id ON public.loan_applicants USING btree (loan_property_id);
CREATE INDEX idx_loan_applicant_borrower_applicant_id ON public.loan_applicants USING btree (borrower_applicant_id);
COMMENT
ON TABLE public.loan_applicants IS 'loan applicant';

-- Column comments
COMMENT
ON COLUMN public.loan_applicants."type" IS 'Individual,Company';
COMMENT
ON COLUMN public.loan_applicants.invitation_status IS 'Accepted,Pending,Rejected';


CREATE TABLE public.loan_borrowers
(
    id                 serial4      NOT NULL,
    created_at         timestamptz(6) NOT NULL DEFAULT now(),
    updated_at         timestamptz(6) NOT NULL DEFAULT now(),
    is_del             int2         NOT NULL DEFAULT 0,
    loan_id            int4         NOT NULL DEFAULT 0,
    "type"             varchar(20)  NOT NULL DEFAULT ''::character varying,
    "name"             varchar(255) NOT NULL DEFAULT ''::character varying,
    tin_type           varchar(10)  NOT NULL DEFAULT 'SSN'::character varying,
    tin                varchar(20)  NOT NULL DEFAULT ''::character varying,
    mailing_address    jsonb        NOT NULL DEFAULT '{}'::jsonb,
    contact_email      varchar(100) NOT NULL DEFAULT ''::character varying,
    contact_phone      varchar(20)  NOT NULL DEFAULT ''::character varying,
    contact_name       varchar(100) NOT NULL DEFAULT ''::character varying,
    borrower_entity_id int4         NOT NULL DEFAULT 0,
    CONSTRAINT loan_borrowers_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_loan_bo_borrower_entity_id ON public.loan_borrowers USING btree (borrower_entity_id);
CREATE INDEX idx_loan_bo_loan_id ON public.loan_borrowers USING btree (loan_id);


CREATE TABLE public.loan_demographics
(
    id                int4        NOT NULL DEFAULT nextval('loan_terms_id_seq'::regclass),
    created_at        timestamptz(6) NOT NULL DEFAULT now(),
    updated_at        timestamptz(6) NOT NULL DEFAULT now(),
    is_del            int2        NOT NULL DEFAULT 0,
    loan_id           int4        NOT NULL DEFAULT 0,
    property          jsonb       NOT NULL DEFAULT '{}'::jsonb,
    ethnicity         jsonb       NOT NULL DEFAULT '{}'::jsonb,
    race              jsonb       NOT NULL DEFAULT '{}'::jsonb,
    sex               varchar(10) NOT NULL DEFAULT 'Male'::character varying,
    extra_demographic jsonb       NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT loan_terms_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_loan_terms_loan_id ON public.loan_demographics USING btree (loan_id);


CREATE TABLE public.loan_documents
(
    id               serial4      NOT NULL,
    loan_id          int4         NOT NULL DEFAULT 0,
    category         varchar(30)  NOT NULL,                               -- Lender,Borrower,Extend
    file_type        varchar(100) NOT NULL DEFAULT ''::character varying,
    status           varchar(30)  NOT NULL DEFAULT ''::character varying, -- Pending Upload,Requested,Ready for Review
    created_at       timestamptz(6) NOT NULL DEFAULT now(),
    updated_at       timestamptz(6) NOT NULL DEFAULT now(),
    is_del           int2         NOT NULL DEFAULT 0,
    "comment"        varchar(255) NOT NULL DEFAULT ''::character varying,
    last_uploaded_at timestamptz(6) NULL,
    CONSTRAINT loan_documents_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_loan_documents_loan_id ON public.loan_documents USING btree (loan_id);

-- Column comments
COMMENT
ON COLUMN public.loan_documents.category IS 'Lender,Borrower,Extend';
COMMENT
ON COLUMN public.loan_documents.status IS 'Pending Upload,Requested,Ready for Review';


CREATE TABLE public.loan_experiences
(
    id                 serial4        NOT NULL,
    created_at         timestamptz(6) NOT NULL DEFAULT now(),
    updated_at         timestamptz(6) NOT NULL DEFAULT now(),
    is_del             int2           NOT NULL DEFAULT 0,
    loan_id            int4           NOT NULL DEFAULT 0,
    ownership_on_title varchar(100)   NOT NULL DEFAULT ''::character varying,
    property_type      varchar(50)    NOT NULL DEFAULT ''::character varying,
    purchase_amount    numeric(10, 2) NOT NULL DEFAULT 0,
    purchase_date      date NULL,
    sold_amount        numeric(10, 2) NOT NULL DEFAULT 0,
    sold_date          date NULL,
    budget_or_cost     numeric(10, 2) NOT NULL DEFAULT 0,
    profit             numeric(10, 2) NOT NULL DEFAULT 0,
    is_tenant_occupied bool           NOT NULL DEFAULT false,
    "type"             varchar(30)    NOT NULL DEFAULT ''::character varying,
    address            jsonb          NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT loan_experiences_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_loan_exp_loan_id ON public.loan_experiences USING btree (loan_id);


CREATE TABLE public.loan_operations
(
    id            serial4     NOT NULL,
    loan_id       int4        NOT NULL DEFAULT 0,
    operator_id   int4        NOT NULL DEFAULT 0,
    operator_name varchar(50) NOT NULL DEFAULT ''::character varying,
    operator_type varchar(50) NOT NULL DEFAULT ''::character varying,
    "action"      varchar(50) NOT NULL DEFAULT ''::character varying,
    "from"        jsonb       NOT NULL DEFAULT '{}'::jsonb,
    "to"          jsonb       NOT NULL DEFAULT '{}'::jsonb,
    metadata      jsonb       NOT NULL DEFAULT '{}'::jsonb,
    created_at    timestamptz(6) NOT NULL DEFAULT now(),
    updated_at    timestamptz(6) NOT NULL DEFAULT now(),
    is_del        int2        NOT NULL DEFAULT 0,
    CONSTRAINT loan_operations_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_loan_operation_loan_id ON public.loan_operations USING btree (loan_id);
CREATE INDEX idx_loan_operation_operator_id ON public.loan_operations USING btree (operator_id);


CREATE TABLE public.loan_properties
(
    id                              serial4        NOT NULL,
    loan_id                         int4           NOT NULL DEFAULT 0,
    borrower_property_id            int4           NOT NULL DEFAULT 0,
    property_type                   varchar(30)    NOT NULL DEFAULT 'Single Family'::character varying, -- Single Family,Townhouse,Multi Family,Condo
    created_at                      timestamptz(6) NOT NULL DEFAULT now(),
    updated_at                      timestamptz(6) NOT NULL DEFAULT now(),
    is_del                          int2           NOT NULL DEFAULT 0,
    living_area_sqft                numeric(8, 2)  NOT NULL DEFAULT 0,
    bedrooms                        int2           NOT NULL DEFAULT 0,
    bathrooms                       int2           NOT NULL DEFAULT 0,
    future_property_type            varchar(30)    NOT NULL DEFAULT ''::character varying,
    future_unit                     int2           NOT NULL DEFAULT 1,
    future_living_area_sqft         numeric(8, 2)  NOT NULL DEFAULT 0,
    future_bedrooms                 int2           NOT NULL DEFAULT 0,
    future_bathrooms                int2           NOT NULL DEFAULT 0,
    "owner"                         varchar(255)   NOT NULL DEFAULT ''::character varying,
    unit                            int2           NOT NULL DEFAULT 1,
    existing_lender_company         varchar(100)   NOT NULL DEFAULT ''::character varying,
    existing_rate                   numeric(4, 2)  NOT NULL DEFAULT 0,
    existing_maturity_date          date NULL,
    purchase_amount                 numeric(10, 2) NOT NULL DEFAULT 0,
    purchase_date                   date NULL,
    has_assignment_fee              bool           NOT NULL DEFAULT false,
    assignment_fee_amount           numeric(10, 2) NOT NULL DEFAULT 0,
    has_subordinate_debt            bool           NOT NULL DEFAULT false,
    subordinate_debt_amount         numeric(10, 2) NOT NULL DEFAULT 0,
    scheduled_close_escrow_date     date NULL,
    liquidity_amount                numeric(10, 2) NOT NULL DEFAULT 0,
    budget                          jsonb          NOT NULL DEFAULT '{}'::jsonb,
    mid_construction                bool           NOT NULL DEFAULT false,
    est_yearly_rental_income_amount numeric(10, 2) NOT NULL DEFAULT 0,
    yearly_liability_expense_amount numeric(10, 2) NOT NULL DEFAULT 0,
    dscr                            numeric(10, 2) NOT NULL DEFAULT 0,
    has_delayed_financing           bool           NOT NULL DEFAULT false,
    has_collateralized_blanket      bool           NOT NULL DEFAULT false,
    has_properties_vacant           bool           NOT NULL DEFAULT false,
    has_short_term_rental           bool           NOT NULL DEFAULT false,
    est_property_value_amount       numeric(10, 2) NOT NULL DEFAULT 0,
    address                         jsonb          NOT NULL DEFAULT '{}'::jsonb,
    arv                             numeric(10, 2) NOT NULL DEFAULT 0,
    CONSTRAINT loan_properties_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_loan_property_borrower_property_id ON public.loan_properties USING btree (borrower_property_id);
CREATE INDEX idx_loan_property_loan_id ON public.loan_properties USING btree (loan_id);

-- Column comments
COMMENT
ON COLUMN public.loan_properties.property_type IS 'Single Family,Townhouse,Multi Family,Condo';


CREATE TABLE public.loan_third_services
(
    id         serial4      NOT NULL,
    created_at timestamptz(6) NOT NULL DEFAULT now(),
    updated_at timestamptz(6) NOT NULL DEFAULT now(),
    is_del     int2         NOT NULL DEFAULT 0,
    loan_id    int4         NOT NULL DEFAULT 0,
    "type"     varchar(20)  NOT NULL DEFAULT ''::character varying,
    "name"     varchar(100) NOT NULL,
    email      varchar(100) NOT NULL DEFAULT ''::character varying,
    phone      varchar(20)  NOT NULL DEFAULT ''::character varying,
    company    varchar(255) NOT NULL DEFAULT ''::character varying,
    CONSTRAINT loan_third_services_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_loan_3rd_svc_loan_id ON public.loan_third_services USING btree (loan_id);


CREATE TABLE public.loans
(
    id                          serial4        NOT NULL,
    borrower_id                 int4           NOT NULL DEFAULT 0,
    lender_company_id           int4           NOT NULL DEFAULT 0,
    stage                       varchar(10)    NOT NULL DEFAULT 'Pending'::character varying,  -- Pending,Active,Archived
    pending_status              varchar(20)    NOT NULL DEFAULT 'Quoting'::character varying,  -- Quoting,Processing,Closing
    loan_purpose                varchar(20)    NOT NULL DEFAULT 'Purchase'::character varying, -- Purchase,Rate-term refi,Cash-out refi
    total_amount                numeric(10, 2) NOT NULL DEFAULT 0,                             -- total loan amount
    initial_amount              numeric(10, 2) NOT NULL DEFAULT 0,                             -- initial loan amount
    closing_date                date NULL,
    holdback_amount             numeric(10, 2) NOT NULL DEFAULT 0,
    rate                        numeric(4, 2)  NOT NULL DEFAULT 0,
    closing_fund                numeric(10, 2) NOT NULL DEFAULT 0,
    created_at                  timestamptz(6) NOT NULL DEFAULT now(),
    updated_at                  timestamptz(6) NOT NULL DEFAULT now(),
    is_del                      int2           NOT NULL DEFAULT 0,
    loan_type                   varchar(255)   NOT NULL DEFAULT 'Bridge'::character varying,   -- Bridge,Rehab,Ground up construction
    paid_off_date               date NULL,                                                     -- paid off date
    principal_amount            numeric(10, 2) NOT NULL DEFAULT 0,                             -- after maturiry date,it is called outstanding amount
    monthly_payment             numeric(10, 2) NOT NULL DEFAULT 0,
    maturity_date               date NULL,
    holdback_released_amount    numeric(10, 2) NOT NULL DEFAULT 0,
    lender_company              varchar(255)   NOT NULL DEFAULT ''::character varying,
    service_company             varchar(255)   NOT NULL DEFAULT ''::character varying,
    business_strategy           _text NULL DEFAULT '{}'::text[],
    "source"                    varchar(20)    NOT NULL DEFAULT 'system'::character varying,   -- borrower,lender,system
    lien_position               varchar(255) NULL DEFAULT ''::character varying,               -- 1st,2nd,3rd
    prepayment_penalty          int2           NOT NULL DEFAULT 0,
    existing_lender_company     varchar(255)   NOT NULL DEFAULT ''::character varying,
    existing_rate               numeric(4, 2)  NOT NULL DEFAULT 0,
    existing_maturity_date      date NULL,
    purchase_amount             numeric(10, 2) NOT NULL DEFAULT 0,
    purchase_date               date NULL,
    has_assignment_fee          bool           NOT NULL DEFAULT false,
    assignment_fee_amount       numeric(10, 2) NOT NULL DEFAULT 0,
    has_subordinate_debt        bool           NOT NULL DEFAULT false,
    subordinate_debt_amount     numeric(10, 2) NOT NULL DEFAULT 0,
    scheduled_close_escrow_date date NULL,
    loan_term                   int2           NOT NULL DEFAULT 0,
    arv                         numeric(10, 2) NOT NULL DEFAULT 0,
    ltv                         numeric(10, 2) NOT NULL DEFAULT 0,
    liquidity_amount            numeric(10, 2) NOT NULL DEFAULT 0,
    budget                      jsonb          NOT NULL DEFAULT '{}'::jsonb,
    mid_construction            bool           NOT NULL DEFAULT false,
    est_credit_score            int2           NOT NULL DEFAULT 0,
    est_yearly_rental_income    numeric(10, 2) NOT NULL DEFAULT 0,
    yearly_liability_expense    numeric(10, 2) NOT NULL DEFAULT 0,
    dscr                        numeric(8, 2)  NOT NULL DEFAULT 0,
    reserves_amount             numeric(10, 2) NOT NULL DEFAULT 0,
    has_delayed_financing       bool           NOT NULL DEFAULT false,
    has_collateralized_blanket  bool           NOT NULL DEFAULT false,
    has_properties_vacant       bool           NOT NULL DEFAULT false,
    has_short_term_rental       bool           NOT NULL DEFAULT false,
    code                        varchar(100)   NOT NULL DEFAULT ''::character varying,
    original_contract_date      date NULL,
    extension_date              date NULL,
    loan_processor              varchar(100)   NOT NULL DEFAULT ''::character varying,
    extend                      jsonb          NOT NULL DEFAULT '{}'::jsonb,
    est_property_value_amount   numeric(10, 2) NOT NULL DEFAULT 0,
    CONSTRAINT idx_loan_code UNIQUE (code),
    CONSTRAINT loans_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_loan_borrower_id ON public.loans USING btree (borrower_id);
CREATE INDEX idx_loan_lender_id ON public.loans USING btree (lender_company_id);

-- Column comments
COMMENT
ON COLUMN public.loans.stage IS 'Pending,Active,Archived';
COMMENT
ON COLUMN public.loans.pending_status IS 'Quoting,Processing,Closing';
COMMENT
ON COLUMN public.loans.loan_purpose IS 'Purchase,Rate-term refi,Cash-out refi';
COMMENT
ON COLUMN public.loans.total_amount IS 'total loan amount';
COMMENT
ON COLUMN public.loans.initial_amount IS 'initial loan amount';
COMMENT
ON COLUMN public.loans.loan_type IS 'Bridge,Rehab,Ground up construction';
COMMENT
ON COLUMN public.loans.paid_off_date IS 'paid off date';
COMMENT
ON COLUMN public.loans.principal_amount IS 'after maturity date,it is called outstanding amount';
COMMENT
ON COLUMN public.loans."source" IS 'borrower,lender,system';
COMMENT
ON COLUMN public.loans.lien_position IS '1st,2nd,3rd';


CREATE TABLE public.users
(
    id            serial4      NOT NULL,
    email         varchar(255) NOT NULL DEFAULT ''::character varying,
    password_hash varchar(255) NOT NULL DEFAULT ''::character varying,
    status        varchar(20) NULL DEFAULT 'unverified'::character varying,    -- active,inactive
    created_at    timestamptz(6) NULL DEFAULT now(),
    updated_at    timestamptz(6) NULL DEFAULT now(),
    "type"        varchar(20)  NOT NULL DEFAULT 'borrower'::character varying, -- borrower,lender
    is_del        int2         NOT NULL DEFAULT 0,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_user_email ON public.users USING btree (email);

-- Column comments
COMMENT
ON COLUMN public.users.status IS 'active,inactive';
COMMENT
ON COLUMN public.users."type" IS 'borrower,lender';