-- create base_templates table
CREATE TABLE "public"."base_templates" (
     "id" serial,
     "created_at" timestamptz NOT NULL DEFAULT now(),
     "updated_at" timestamptz NOT NULL DEFAULT now(),
     "is_del" int2 NOT NULL DEFAULT 0,
     "name" varchar(100) NOT NULL DEFAULT '',
     "type" varchar(20) NOT NULL DEFAULT '',
     "category" varchar(20) NOT NULL DEFAULT '',
     "purpose" varchar(100) NOT NULL DEFAULT '',
     "content" text NOT NULL DEFAULT '',
     "variable_group_schema" jsonb NOT NULL DEFAULT '{}',
     "status" varchar(20) NOT NULL DEFAULT 'Active',
     "description" varchar(255) NOT NULL DEFAULT '',
     PRIMARY KEY ("id")
);

CREATE INDEX "idx_base_tpl_name" ON "public"."base_templates" ("name");
CREATE INDEX "idx_base_tpl_purpose" ON "public"."base_templates" ("purpose");
COMMENT ON COLUMN "public"."base_templates"."type" IS 'Borrower/Lender';
COMMENT ON COLUMN "public"."base_templates"."category" IS 'Email/Contract';
COMMENT ON COLUMN "public"."base_templates"."status" IS 'Active/Disabled';

-- create organization_templates table
CREATE TABLE "public"."organization_templates" (
   "id" serial,
   "created_at" timestamptz NOT NULL DEFAULT now(),
   "updated_at" timestamptz NOT NULL DEFAULT now(),
   "is_del" int2 NOT NULL DEFAULT 0,
   "organization_id" int4 NOT NULL DEFAULT 0,
   "name" varchar(100) NOT NULL DEFAULT '',
   "type" varchar(20) NOT NULL DEFAULT '',
   "category" varchar(20) NOT NULL DEFAULT '',
   "purpose" varchar(100) NOT NULL DEFAULT '',
   "content" text NOT NULL DEFAULT '',
   "variable_group_schema" jsonb NOT NULL DEFAULT '{}',
   "status" varchar(20) NOT NULL DEFAULT 'Active',
   "description" varchar(255) NOT NULL DEFAULT '',
   PRIMARY KEY ("id")
);

CREATE INDEX "idx_org_tpl_name" ON "public"."organization_templates" ("name");
CREATE INDEX "idx_org_tpl_purpose" ON "public"."organization_templates" ("purpose");
CREATE INDEX "idx_org_tpl_org_id" ON "public"."organization_templates" ("organization_id");
COMMENT ON COLUMN "public"."organization_templates"."type" IS 'Borrower/Lender';
COMMENT ON COLUMN "public"."organization_templates"."category" IS 'Email/Contract';
COMMENT ON COLUMN "public"."organization_templates"."status" IS 'Active/Disabled';

ALTER TABLE "public"."organization_templates"
    ADD COLUMN "base_template_id" int4 NOT NULL DEFAULT 0;
CREATE INDEX "idx_org_base_tpl_id" ON "public"."organization_templates" ("base_template_id");

ALTER TABLE "public"."base_templates"
    ADD COLUMN "example_variables" jsonb NOT NULL DEFAULT '{}';

ALTER TABLE "public"."organization_templates"
    ADD COLUMN "example_variables" jsonb NOT NULL DEFAULT '{}';
ALTER TABLE "public"."base_templates"
    ADD COLUMN "scope" varchar(20) NOT NULL DEFAULT 'System';
COMMENT
ON COLUMN "public"."base_templates"."scope" IS 'System/Organization';

INSERT INTO "public"."base_templates" ("name", "type", "category", "purpose", "content", "variable_group_schema",
                                       "status", "description", "example_variables", "scope")
VALUES ('preApproved', 'Lender', 'Contract', 'preApprovedLetter', '<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Pre-Approval Letter Generator</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Great+Vibes&display=swap" rel="stylesheet">
    <style>
        * {
            box-sizing: border-box;
        }
        .page {
            width: 100%;
            height: 100%;
            /* min-height: 1000px; */
            background-color: #ffffff;
            box-shadow: 0 0 15px rgba(0, 0, 0, 0.1);
            padding: 30px 50px 30px 50px;
            position: relative;
            display: flex;
            flex-direction: column;
            overflow: hidden;
            overflow-y: auto;
            font-family: ''Segoe UI'', Tahoma, Geneva, Verdana, sans-serif;
        }
        /* Letter Styles */
        .header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 5px;
        }
        .logo {
            width: 120px;
            height: 100px;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .logo svg {
            width: 100%;
            height: 100%;
        }
        .company-info {
            text-align: right;
            font-size: 12px;
            color: #1a3c8b;
            line-height: 1.4;
        }
        .company-info .name {
            font-size: 28px;
            font-weight: 500;
            margin-bottom: 2px;
            color: #1a3c8b;
            letter-spacing: -0.5px;
        }
        .company-info .slogan {
            color: #3b6fd6;
            margin-bottom: 8px;
            font-style: italic;
        }
        /* Layout for Title + Officer Info */
        .content-row {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-top: 10px;
            margin-bottom: 10px;
        }
        .content-left {
            flex: 1;
        }
        .content-right {
            text-align: left;
            margin-left: 20px;
        }
        .officer-info {
            font-size: 13px;
            line-height: 1.4;
            color: #1a3c8b;
        }
        .officer-name {
            font-weight: 800;
            font-size: 18px;
            margin-bottom: 4px;
            font-family: ''Times New Roman'', serif;
            letter-spacing: 0.5px;
        }
        .title {
            margin-top: 0;
            margin-bottom: 10px;
            color: #1a3c8b;
            font-weight: bold;
            font-size: 26px;
            line-height: 1.2;
        }
        .title span {
            display: block;
        }
        .meta {
            font-size: 13px;
            margin-bottom: 15px;
        }
        .meta .date {
            font-weight: bold;
            text-decoration: underline;
        }
        .body-text {
            font-size: 13px;
            line-height: 1.6;
            margin-bottom: 10px;
            color: #333;
        }
        .terms-table {
            width: 100%;
            border-collapse: collapse;
            margin: 15px 0;
            font-size: 13px;
        }
        .terms-table th,
        .terms-table td {
            border: 1px solid #333;
            padding: 4px 8px;
            vertical-align: middle;
        }
        .terms-table th {
            text-align: left;
            background-color: #f8f9fa;
            /* Light gray */
            width: 40%;
            font-weight: 600;
        }
        .footer-list {
            font-size: 13px;
            line-height: 1.6;
            margin-left: 0;
            padding-left: 20px;
            color: #333;
            margin-bottom: 10px;
            margin-top: 5px;
        }
        .footer-list li {
            margin-bottom: 2px;
        }
        /* Signature and Footer */
        .signature-section {
            margin-top: 10px;
            margin-bottom: 10px;
        }
        .signature-cursive {
            font-family: ''Great Vibes'', cursive;
            font-size: 30px;
            font-weight: normal;
            margin-bottom: 5px;
        }
        .signature-print {
            font-family: ''Times New Roman'', serif;
            font-size: 14px;
            font-style: italic;
        }
        .page-footer {
            margin-top: 10px;
            border-top: 2px solid #1a3c8b;
            padding-top: 10px;
            display: flex;
            justify-content: space-between;
            color: #1a3c8b;
            font-style: italic;
            font-size: 12px;
        }
        @media print {
            .editor-panel {
                display: none;
            }
            .page {
                box-shadow: none;
                margin: 0;
                width: 100%;
            }
        }
    </style>
</head>
<body>
  <!-- <div style="overflow: hidden;padding:10px;"> -->
        <div class="page">
            <div class="header">
                <!-- Logo -->
                <div class="logo" id="preview-logo">
                    <!-- Using img with data URI for better PDF generation compatibility -->
                    <img src=''[[ .system.organization.Logo ]]'' alt="Logo" style="width: 100%; height: 100%;">
                </div>
                <!-- Company Info -->
                <div class="company-info">
                    <div class="name" id="preview-company-name">[[ .system.organization.Name ]]</div>
                    <div class="slogan" id="preview-slogan">[[ .system.organization.Tagline ]]</div>
                </div>
            </div>
            <div class="content-row">
                <div class="content-left">
                    <div class="title">
                        <span id="preview-title-1" data-variable=".letterTitle1">{{ .letterTitle1 }}</span>
                        <span id="preview-title-2" data-variable=".letterTitle2">{{ .letterTitle2 }}</span>
                    </div>
                    <div class="meta">
                        <span class="date" id="preview-date">{{ .pdfCreateTime }}</span>
                    </div>
                    <div class="body-text" style="margin-top: 10px;">
                        Dear <span id="preview-borrower" style="font-weight:bold;">{{ .loanBorrower }}</span>:
                    </div>
                </div>

                <div class="content-right">
                    <div class="officer-info">
                        <div class="officer-name" id="preview-officer-name"  data-variable=".officerName">{{ .officerName }}</div>
                        <div id="preview-officer-title" data-variable=".officerTitle">{{ .officerTitle }}</div>
                        <div>Mortgage Consultant</div>
                        <div>NMLS: <span id="preview-nmls" data-variable=".nmls">{{ .nmls }}</span></div>
                        <div>E-Mail: <span id="preview-email" data-variable=".email">{{ .email }}</span></div>
                        <div>Direct: <span id="preview-phone" data-variable=".phone">{{ .phone }}</span></div>
                    </div>
                </div>
            </div>

            <div class="body-text" id="preview-body-1" data-variable=".openingParagraph">
                <!-- <span id="preview-company-name-text">AI Mortgage Lending</span>  -->
                {{ .openingParagraph }}
            </div>

            <table class="terms-table">
                <tr>
                    <th>Sales Price</th>
                    <td id="preview-sales-price" data-variable=".salesPrice">{{ .salesPrice }}</td>
                </tr>
                <tr>
                    <th>Max Loan Amount</th>
                    <td id="preview-max-loan-amount" data-variable=".maxLoanAmount">{{ .maxLoanAmount }}</td>
                </tr>
                <!-- <tr>
                    <th>Down Payment</th>
                    <td id="preview-down-payment" data-variable=".downPayment">{{ .downPayment }}</td>
                </tr> -->
                <!-- <tr>
                    <th>Loan Program</th>
                    <td id="preview-loan-program" data-variable=".loanProgram">{{ .loanProgram }} Yrs Amortization</td>
                </tr> -->
                <tr>
                    <th>Terms</th>
                    <td id="preview-terms" data-variable=".terms">{{ .terms }} Months</td>
                </tr>
                <tr>
                    <th>Interest Rate</th>
                    <td id="preview-interest-rate"  data-variable=".interestRate">{{ .interestRate }} %</td>
                </tr>
            </table>

            <div class="body-text" id="preview-body-2" style="white-space: pre-wrap;" data-variable=".closingParagraph">{{ .closingParagraph }}</div>

            <div class="signature-section">
                <div style="margin-bottom: 15px;">Sincerely,</div>
                <div class="signature-cursive" id="preview-signature-cursive"  data-variable=".officerName">{{ .officerName }}</div>
                <div class="signature-print" id="preview-signature-print"  data-variable=".officerName">{{ .officerName }}</div>
            </div>

            <div class="page-footer">
                <div>Address: <span id="preview-company-address" data-variable=".companyAddress">{{ .companyAddress }}</span>
                </div>
                <div>Phone: <span id="preview-company-phone" data-variable=".companyPhone">{{ .companyPhone }}</span></div>
                <div>Fax: <span id="preview-company-fax" data-variable=".companyFax">{{ .companyFax }}</span></div>
            </div>
        </div>
  <!-- </div> -->
</body>

</html>', '[
  {
    "name": "Header",
    "variables": [
      {
        "key": ".letterTitle1",
        "type": "String",
        "label": "Letter title 1",
        "format": "Input",
        "source": "Template",
        "fixed_value": "HOME BUYER",
        "is_editable": true
      },
      {
        "key": ".letterTitle2",
        "type": "String",
        "label": "Letter title 2",
        "format": "Input",
        "source": "Template",
        "fixed_value": "LOAN PRE-APPROVAL",
        "is_editable": true
      },
      {
        "key": ".officerName",
        "type": "String",
        "label": "Officer name",
        "format": "Input",
        "source": "Template",
        "fixed_value": "Angel Huang",
        "is_editable": true
      },
      {
        "key": ".officerTitle",
        "type": "String",
        "label": "Officer title",
        "format": "Input",
        "source": "Template",
        "fixed_value": "Senior Loan Officer",
        "is_editable": true
      },
      {
        "key": ".nmls",
        "type": "String",
        "label": "NMLS ID",
        "format": "Input",
        "source": "Template",
        "fixed_value": "2651574",
        "is_editable": true
      },
      {
        "key": ".email",
        "type": "String",
        "label": "Email",
        "format": "Email",
        "source": "Template",
        "fixed_value": "rebyangel@gmail.com",
        "is_editable": true
      },
      {
        "key": ".phone",
        "type": "String",
        "label": "Phone",
        "format": "Phone",
        "source": "Template",
        "fixed_value": "916-276-3242",
        "is_editable": true
      },
      {
        "key": ".pdfCreateTime",
        "type": "String",
        "label": "Pdf create time",
        "format": "Input",
        "source": "Dynamic",
        "fixed_value": "1",
        "is_editable": false
      },
      {
        "key": ".loanBorrower",
        "type": "String",
        "label": "Loan Borrower",
        "format": "Input",
        "source": "Dynamic",
        "fixed_value": "1",
        "is_editable": false
      }
    ]
  },
  {
    "name": "Opening Paragraph",
    "variables": [
      {
        "key": ".openingParagraph",
        "type": "String",
        "label": "Paragraph",
        "format": "TextArea",
        "source": "Template",
        "fixed_value": "22",
        "is_editable": true
      }
    ]
  },
  {
    "name": "Quote Detail",
    "variables": [
      {
        "key": ".salesPrice",
        "type": "String",
        "label": "Sales Price",
        "format": "Currency",
        "source": "Dynamic",
        "fixed_value": "1,200,000.02",
        "is_editable": true
      },
      {
        "key": ".maxLoanAmount",
        "type": "String",
        "label": "Max Loan Amount",
        "format": "Currency",
        "source": "Dynamic",
        "fixed_value": "840,000",
        "is_editable": true
      },
      {
        "key": ".terms",
        "type": "String",
        "label": "Terms(Months)",
        "format": "Input",
        "source": "Dynamic",
        "fixed_value": "6",
        "is_editable": true
      },
      {
        "key": ".interestRate",
        "type": "String",
        "label": "Interest Rate(%)",
        "format": "Input",
        "source": "Dynamic",
        "fixed_value": "60",
        "is_editable": true
      }
    ]
  },
  {
    "name": "Closing Paragraph",
    "variables": [
      {
        "key": ".closingParagraph",
        "type": "String",
        "label": "Paragraph",
        "format": "TextArea",
        "source": "Template",
        "fixed_value": "1",
        "is_editable": true
      }
    ]
  },
  {
    "name": "Footer",
    "variables": [
      {
        "key": ".companyAddress",
        "type": "String",
        "label": "Company Address",
        "format": "Input",
        "source": "Template",
        "fixed_value": "388 E Valley Blvd suite 200, Alhambra, CA 91801",
        "is_editable": true
      },
      {
        "key": ".companyPhone",
        "type": "String",
        "label": "Company Phone",
        "format": "Phone",
        "source": "Template",
        "fixed_value": "626-588-8688",
        "is_editable": true
      },
      {
        "key": ".companyFax",
        "type": "String",
        "label": "Company Fax",
        "format": "Phone",
        "source": "Template",
        "fixed_value": "626-737-6035",
        "is_editable": true
      }
    ]
  }
]', 'Active', 'dddd', '{
  "salesPrice": 8888,
  "maxLoanAmount": 999
}', 'Organization');

INSERT INTO "public"."base_templates" ("name", "type", "category", "purpose", "content", "variable_group_schema",
                                       "status", "description", "example_variables", "scope")
VALUES ('preApprovedLetterEmail', 'Lender', 'Email', 'preApprovedLetterEmailTemplate', '<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <title>Pre-Approval Letter Generator</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;500;700&display=swap" rel="stylesheet">
  <meta name="viewport" content="width=device-width, initial-scale=1">
</head>

<body>
  <div style="width:600px;margin:auto;font-family:Roboto;">
    <div style="font-size:16px;padding-bottom:16px;">Dear {{ .borrowerName }},</div>
    <div style="font-size:16px;">Attached please find your
      <span style="font-weight:500;">Loan Pre-Approval Letter</span> , which outlines the loan amount, terms, and
      conditions based on the information reviewed to date.</div>
    <div style="font-size:14px;padding:24px 0 16px;">Please review the document carefully. If you have any questions,
      need revisions, or are ready to move forward with the next steps, feel free to contact me at any time. I would be
      happy to assist you.</div>
    <div style="font-size:14px;">Thank you for choosing
      <span style="font-weight:500;">[[ .system.organization.Name ]]</span></div>
  </div>
</body>

</html>', '[
  {
    "name": "BorrowerInfo",
    "variables": [
      {
        "key": ".borrowerName",
        "type": "String",
        "label": "Borrower name",
        "format": "Input",
        "source": "Dynamic",
        "fixed_value": "",
        "is_editable": true
      }
    ]
  }
]', 'Active', '用于preApproved提上来的loan，发送preApprovedLetter附件时的邮件模板', '{}', 'Organization');

