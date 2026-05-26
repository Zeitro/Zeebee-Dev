# Borrower Property Transaction 技术设计方案

| 字段 | 值 |
|------|-----|
| 项目 | Zeebee Backend |
| 模块 | Borrower Property Transaction |
| 分支 | `001-borrower-property-transaction` |
| 作者 | Wei Zeng |
| 创建日期 | 2026-05-19 |
| 最后更新 | 2026-05-26 |
| 版本 | v3 |
| 状态 | 开发中（Phase 1-13 已完成，Phase 14 验证待完成） |

---

## 1. 背景与目标

### 1.1 业务背景

Zeebee 平台为贷款经纪人（Loan Officer）管理借款人（Borrower）的资产（Property）和贷款（Loan）流程。当前系统支持 Property 和 Loan 的直接关联，但缺少对交易过程（Transaction）的管理。

实际业务中，贷款申请前需要为 Property 准备大量交易文件（合同、托管、保险、过户包等），这些文件的准备进度直接影响贷款审批。需要一个独立的 Transaction 模块来管理交易生命周期和相关文档。

### 1.2 目标

1. **Transaction 管理** — 为 Property 添加 Transaction 实体，支持完整生命周期管理
2. **通用文档系统** — 基于模板的文档树管理，支持文件上传/删除，可复用于其他业务
3. **Loan 关联** — Loan 可通过 Transaction 间接关联 Property
4. **操作审计** — 统一操作日志，记录所有业务变更

### 1.3 范围边界

**本期实现**:
- Transaction CRUD + 状态切换
- 基于 S3 的文档管理（Manual 上传）
- Checklist 自动统计
- Loan 关联 Transaction
- 操作日志记录与查询
- Transaction 分页搜索

**不在本期**:
- Google Drive 文件集成（数据模型已预留，逻辑未实现）
- Document Template CRUD 管理（手动配置入库）
- 操作日志的实时推送/通知

---

## 2. 模块功能要点

### 2.1 Transaction 管理

| 操作 | 说明 | 约束 |
|------|------|------|
| Create | 选择 Property + Type，自动生成 code 和文档树 | 同一 Property 可创建多个同类型 Transaction |
| Detail | 根据 ID 获取详情 | — |
| Update | 编辑业务字段 | 仅限: coe_date, purchase_amount, emd_amount, property_type, pm, strategy, note |
| Delete | 软删除 | 仅未绑定 Loan 的 Transaction 可删除，递归删除关联 Document 和 File |
| SwitchStatus | Pending/Closed/Cancelled 任意切换 | 无前置状态约束 |
| All | 返回下拉选项列表 | 用于 Loan 关联选择 |
| Page | 分页搜索 | 支持 keyword(地址模糊匹配) + type + status 筛选 |

**Transaction Type**: Purchase, Sale, Re-finance, Change Title, Assignment

**Transaction Status**: Pending, Closed, Cancelled

**Code 规则**: `TC` + 年份2位 + base62(毫秒时间戳)，例: `TC26ttdhd10`

### 2.2 Document 管理（通用文档系统）

| 操作 | 说明 | 约束 |
|------|------|------|
| Create Document | 创建叶子节点同级 document | 同 level 内 name 不重复，禁止创建非叶子节点 |
| Upload File | 上传文件到叶子 document | 仅 Manual 类型，storage_type=Aws S3 |
| Delete File | 删除文件 | 仅限 upload_type=Manual 的文件 |
| Document Tree | 获取完整文档树 | 递归返回，叶子节点含文件列表 |
| Leaf List | 获取叶子节点列表 | 含每节点文件数量 |

**模板机制**: 创建 Transaction 时从 `document_templates` 按 type 加载 structure，合并 common + per-type 模板后递归创建文档目录树。

### 2.3 Checklist 管理

- 自动统计 required/optional 分组的文档上传进度
- 每次文件操作（上传/删除/创建 document）自动重算
- `all_checked` 为手动确认字段，计算时保留原值

### 2.4 Loan 关联

- Apply Loan: property 与 transaction 二选一（XOR），同时提供以transaction为准
- Draft Loan: 灵活策略，可只提供其中一个或都不提供

### 2.5 操作日志

- 统一 `operations` 表，通过 `entity_type` + `entity_id` 多态关联
- 支持 Transaction 和 Loan 两种实体类型
- 提供 Stats（统计）和 List（列表）两种查询接口

---

## 3. 技术架构

### 3.1 系统架构

本模块遵循项目整体分层架构:

```
Client (Web/Mobile)
    │ HTTP
    ▼
┌─────────────────────────────────────────────┐
│ API Layer — borrower-api                     │
│  app/borrower/api/internal/logic/            │
│  (参数验证、请求透传、响应整合)                  │
└──────────────┬──────────────────────────────┘
               │ gRPC
               ▼
┌─────────────────────────────────────────────┐
│ RPC Layer — borrower-rpc                     │
│  app/borrower/rpc/internal/logic/            │
│  (核心业务逻辑、事务管理)                       │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────┐
│ Shared Domain Layer — domains/zeebee/        │
│  model → repo → logic                        │
│  (数据模型、仓储、领域逻辑)                      │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────┐
│ Infrastructure — PostgreSQL / Redis / S3     │
└─────────────────────────────────────────────┘
```

### 3.2 模块依赖关系

```
borrower-api
  ├── borrower_property_transaction.api (新增)
  ├── document.api (新增)
  └── borrower_property.api (修改: 新增 transaction 分页查询)

borrower-rpc
  ├── borrowerpropertytransaction/ (新增: 10个 RPC Logic)
  ├── document/ (新增: 5个 RPC Logic)
  └── borrowerpropertytransaction/*.proto (新增)

domains/zeebee
  ├── model/
  │   ├── borrower_property_transaction.go (新增)
  │   ├── document.go (新增)
  │   ├── entity_document.go (新增)
  │   ├── operation.go (新增)
  │   ├── document_template.go (修改: DocumentTemplateStructureMap)
  │   └── file.go, loan.go (修改: 扩展字段)
  ├── repo/
  │   ├── borrower_property_transaction.go (新增)
  │   ├── document.go (新增)
  │   ├── entity_document.go (新增)
  │   └── operation.go (新增)
  ├── data/dto/
  │   └── borrower_property_transaction.go (新增: 含 Page/Checklist DTO)
  ├── constants/
  │   ├── borrower_property_transaction.go (新增)
  │   └── operation.go (新增)
  └── logic/
      ├── checklist_logic.go (新增)
      └── operation_logic.go (新增)
```

### 3.3 技术选型

| 组件 | 技术 | 说明 |
|------|------|------|
| Web 框架 | go-zero | API Gateway + RPC Server |
| ORM | GORM | 数据库操作 |
| 数据库 | PostgreSQL | 关系型存储 |
| 文件存储 | AWS S3 | Presigned URL 直传 |
| 缓存 | Redis | Token/Session |
| 数据转换 | copierx.Trans | DTO↔Model↔BO↔Proto |
| ID 生成 | Serial (PG) | 自增主键 |
| Code 生成 | base62 | Transaction 唯一标识 |

---

## 4. 流程设计

### 4.1 创建 Transaction 流程

```
API: POST /borrower_property_transaction
         │
         ▼
    ┌─────────────┐
    │ 验证 Property │  ← 检查存在性 + 组织归属
    └──────┬──────┘
           │
           ▼
    ┌─────────────────┐
    │ 生成 Code        │  ← TC + 年份2位 + base62(毫秒)
    └──────┬──────────┘
           │
           ▼
    ┌─────────────────┐
    │ 创建 Transaction │  ← status=Pending
    └──────┬──────────┘
           │
           ▼
    ┌──────────────────────┐
    │ 加载 DocumentTemplate │  ← common + per-type 合并
    └──────┬───────────────┘
           │
           ▼
    ┌───────────────────┐
    │ 递归创建文档目录树    │  ← 设置 metadata.is_required
    └──────┬────────────┘
           │
           ▼
    ┌──────────────────┐
    │ 创建 EntityDocument│  ← 绑定 transaction ↔ root document
    └──────┬───────────┘
           │
           ▼
    ┌────────────────────┐
    │ 初始化 Checklist     │  ← RecalculateChecklist
    └──────┬─────────────┘
           │
           ▼
    ┌──────────────┐
    │ 记录操作日志   │  ← OperationLogic
    └──────────────┘
```

### 4.2 文件上传流程

前端直传 S3，后端仅负责生成 Presigned URL 和记录文件信息。

```
Step 1: 前端请求上传 URL
API: POST /borrower_property_transaction/:tid/document/:did/file
         │
         ▼
    ┌───────────────────┐
    │ 验证 Document 归属  │  ← 属于该 Transaction 的文档树
    └──────┬────────────┘
           │
           ▼
    ┌─────────────┐
    │ 验证叶子节点  │  ← 只能上传到叶子 Document
    └──────┬──────┘
           │
           ▼
    ┌──────────────────────────┐
    │ 生成 S3 Presigned PUT URL │  ← 返回给前端
    └──────────────────────────┘

Step 2: 前端直传文件到 S3（后端不参与）

Step 3: 前端确认上传（同一 API 的回调/确认机制）
         │
         ▼
    ┌──────────────────┐
    │ 创建 File 记录     │  ← storage_type=Aws S3, upload_type=Manual
    └──────┬───────────┘
           │
           ▼
    ┌──────────────────────────┐
    │ 更新 Document last_uploaded_at │  ← 同时更新所有 parent
    └──────┬───────────────────┘
           │
           ▼
    ┌─────────────────────┐
    │ RecalculateChecklist │  ← 重算统计
    └──────┬──────────────┘
           │
           ▼
    ┌──────────────┐
    │ 记录操作日志   │
    └──────────────┘
```

### 4.3 Checklist 计算流程

```
RecalculateChecklist(transactionId)
         │
         ▼
    ┌─────────────────────┐
    │ 获取 EntityDocument   │  ← entity_type + entity_id
    └──────┬──────────────┘
           │
           ▼
    ┌───────────────────┐
    │ 获取根 Document     │
    └──────┬────────────┘
           │
           ▼
    ┌───────────────────────┐
    │ 递归遍历文档树          │
    │ 按 metadata.is_required │
    │ 分组统计:               │
    │ - required: total, updated │
    │ - optional: total, updated │
    └──────┬────────────────┘
           │
           ▼
    ┌───────────────────────┐
    │ 合并原 all_checked 值   │  ← 保留手动确认状态
    └──────┬────────────────┘
           │
           ▼
    ┌───────────────────────┐
    │ 写入 transaction.checklist │
    └───────────────────────┘
```

### 4.4 Loan Apply 二选一流程

```
ApplyLoan(propertyId, transactionId)
         │
         ▼
    ┌────────────────────────────┐
    │ propertyId 和 transactionId │
    │ 是否恰好只有一个有值?         │
    └──────┬──────────┬─────────┘
           │          │
        是 │          │ 否
           │          │
           ▼          ▼
    ┌──────────┐  ┌────────────┐
    │ 继续处理   │  │ 返回参数错误 │
    └──────────┘  └────────────┘
```

---

## 5. 数据库设计

### 5.1 Entity 关系图

```
BorrowerProperty 1───* BorrowerPropertyTransaction
                              │
                              ├── 1:1 EntityDocument (entity_type='BorrowerPropertyTransaction')
                              │         │
                              │         └── 1:1 Document (root, type=Root)
                              │                   │
                              │                   └── * Document (children, type=Branch/Leaf)
                              │                             │
                              │                             └── * File (entity_type='Document')
                              │
                              ├── checklist (cached JSONB, write-on-update)
                              └── * Operation (entity_type='BorrowerPropertyTransaction')

Loan ─── 0..1 BorrowerPropertyTransaction
  │
  └── * Operation (entity_type='Loan')

DocumentTemplate (static reference, no FK relationship)
  structure: DocumentTemplateStructureMap (map[string][]DocumentTemplateStructure)
```

### 5.2 新增表

#### borrower_property_transactions

```sql
CREATE TABLE "public"."borrower_property_transactions" (
  "id" serial,
  "borrower_property_id" int4 NOT NULL DEFAULT 0,
  "organization_id" int4 NOT NULL DEFAULT 0,
  "code" varchar(20) NOT NULL DEFAULT '',
  "type" varchar(20) NOT NULL DEFAULT 'Purchase',
  "status" varchar(20) NOT NULL DEFAULT 'Pending',
  "coe_date" date,
  "purchase_amount" decimal(10,2) NOT NULL DEFAULT 0,
  "emd_amount" decimal(10,2) NOT NULL DEFAULT 0,
  "property_type" varchar(50) NOT NULL DEFAULT '',
  "pm" varchar(255) NOT NULL DEFAULT '',
  "strategy" varchar(100) NOT NULL DEFAULT '',
  "note" varchar(255) NOT NULL DEFAULT '',
  "checklist" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "is_del" int2 NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_bp_trans_property_id" ON "public"."borrower_property_transactions" ("borrower_property_id");
CREATE INDEX "idx_bp_trans_org_id" ON "public"."borrower_property_transactions" ("organization_id");
CREATE UNIQUE INDEX "idx_bp_trans_code" ON "public"."borrower_property_transactions" ("code");
```

#### operations

```sql
CREATE TABLE "public"."operations" (
  "id" serial,
  "entity_id" int4 NOT NULL DEFAULT 0,
  "entity_type" varchar(100) NOT NULL DEFAULT '',
  "operator_id" int4 NOT NULL DEFAULT 0,
  "operator_name" varchar(50) NOT NULL DEFAULT '',
  "operator_type" varchar(50) NOT NULL DEFAULT '',
  "operation" varchar(50) NOT NULL DEFAULT '',
  "note" varchar(255) NOT NULL DEFAULT '',
  "metadata" jsonb NOT NULL DEFAULT '{}',
  "organization_id" int4 NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "is_del" int2 NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_operation_entity" ON "public"."operations" ("entity_id", "entity_type");
CREATE INDEX "idx_operation_org_id" ON "public"."operations" ("organization_id");
```

#### entity_documents

```sql
CREATE TABLE "public"."entity_documents" (
  "id" serial,
  "entity_type" varchar(20) NOT NULL DEFAULT '',
  "entity_id" int4 NOT NULL DEFAULT 0,
  "root_document_id" int4 NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "is_del" int2 NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_entity_doc_entity" ON "public"."entity_documents" ("entity_type", "entity_id");
CREATE INDEX "idx_entity_doc_root_id" ON "public"."entity_documents" ("root_document_id");
```

#### documents

```sql
CREATE TABLE "public"."documents" (
  "id" serial,
  "parent_id" int4 NOT NULL DEFAULT 0,
  "type" varchar(20) NOT NULL DEFAULT '',
  "name" varchar(255) NOT NULL DEFAULT '',
  "level" int4 NOT NULL DEFAULT 0,
  "sort_order" int4 NOT NULL DEFAULT 0,
  "last_uploaded_at" timestamptz,
  "metadata" jsonb NOT NULL DEFAULT '{}',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "is_del" int2 NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);

CREATE INDEX "idx_doc_pid" ON "public"."documents" ("parent_id");
CREATE INDEX "idx_doc_name" ON "public"."documents" ("name");
```

#### document_templates

```sql
CREATE TABLE "public"."document_templates" (
  "id" serial,
  "name" varchar(255) NOT NULL DEFAULT '',
  "type" varchar(20) NOT NULL DEFAULT '',
  "structure" json NOT NULL DEFAULT '{}',
  "status" varchar(20) NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "is_del" int2 NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
```

### 5.3 修改的表

```sql
-- files 表扩展
ALTER TABLE "public"."files"
  ADD COLUMN "file_name" varchar(255) NOT NULL DEFAULT '',
  ADD COLUMN "storage_type" varchar(20) NOT NULL DEFAULT 'Aws S3',
  ADD COLUMN "external_id" varchar(255) NOT NULL DEFAULT '',
  ADD COLUMN "upload_type" varchar(20) NOT NULL DEFAULT 'Manual';
CREATE INDEX "idx_files_external_id" ON "public"."files" ("external_id");

-- loans 表扩展
ALTER TABLE public.loans
  ADD COLUMN borrower_property_transaction_id int4 DEFAULT 0 NOT NULL;
CREATE INDEX idx_loan_borrower_property_trans_id ON public.loans (borrower_property_transaction_id);
```

### 5.4 状态机

```
                    ┌──────────┐
          ┌─────────│ Pending  │─────────┐
          │         └──────────┘         │
          │              │               │
          ▼              │               ▼
    ┌──────────┐         │         ┌──────────┐
    │ Cancelled│◄────────┘────────►│  Closed  │
    └──────────┘                    └──────────┘
          ▲                              │
          └──────────────────────────────┘
```

三状态（Pending/Closed/Cancelled）可任意双向切换，无前置约束。

### 5.5 Checklist JSONB 结构

```json
{"all_checked": true, "optional_documents": [{"name": "1031", "children": [{"name": "1031 Exchange Documents", "document_id": 301, "is_required": false, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 300, "is_required": false, "sub_document_count": 1, "uploaded_document_count": 0}, {"name": "Lease", "children": [{"name": "Lease Agreement", "document_id": 303, "is_required": false, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "TENANT ESTOPPEL CERTIFICATE", "document_id": 304, "is_required": false, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 302, "is_required": false, "sub_document_count": 2, "uploaded_document_count": 0}, {"name": "HOA", "children": [{"name": "See \"HOA\" Sheet", "document_id": 306, "is_required": false, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 305, "is_required": false, "sub_document_count": 1, "uploaded_document_count": 0}, {"name": "Solar", "children": [{"name": "See \"Solar\" Sheet", "document_id": 308, "is_required": false, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 307, "is_required": false, "sub_document_count": 1, "uploaded_document_count": 0}], "required_documents": [{"name": "Contract", "children": [{"name": "Original RPA", "document_id": 267, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Assignment Agreement AOAA", "document_id": 268, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Counter Offer", "document_id": 269, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Addendum (Extension)", "document_id": 270, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Seller financing Docs", "document_id": 271, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "JV Agreement", "document_id": 272, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 266, "is_required": true, "sub_document_count": 6, "uploaded_document_count": 0}, {"name": "Escrow", "children": [{"name": "Wire Instructions", "document_id": 274, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "EMD Payment Proof", "document_id": 275, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "EMD Receipt", "document_id": 276, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Prelim Report", "document_id": 277, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "NHD Report", "document_id": 278, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Colored map", "document_id": 279, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "FinCen", "document_id": 280, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Buyer Open Package", "document_id": 281, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "LOI", "document_id": 282, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 273, "is_required": true, "sub_document_count": 9, "uploaded_document_count": 0}, {"name": "Disclosures", "document_id": 283, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Loan", "children": [{"name": "Application Form", "document_id": 285, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Loan Document", "document_id": 286, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Contact form", "document_id": 287, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Appraisal Report (Appraisal Payment confirmation and receipt)", "document_id": 288, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 284, "is_required": true, "sub_document_count": 4, "uploaded_document_count": 0}, {"name": "Insurance", "children": [{"name": "EOI", "document_id": 290, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Payment Receipt", "document_id": 291, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Full Policy", "document_id": 292, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Flood Insurance", "document_id": 293, "is_required": false, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 289, "is_required": true, "sub_document_count": 4, "uploaded_document_count": 0}, {"name": "Closing Package", "children": [{"name": "Buyer Estimated Statement", "document_id": 295, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Final Buyer Statement", "document_id": 296, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "All Signed Documents", "document_id": 297, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Recorded GD", "document_id": 298, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}, {"name": "Recorded DOT", "document_id": 299, "is_required": true, "sub_document_count": 0, "uploaded_document_count": 0}], "document_id": 294, "is_required": true, "sub_document_count": 5, "uploaded_document_count": 0}], "total_document_count": 34, "total_optional_document_count": 5, "total_required_document_count": 29, "total_optional_updated_document_count": 0, "total_required_updated_document_count": 0}
```

### 5.6 Document metadata JSONB 结构

```json
{
  "is_required": true,
  "status": "Pending",
  "create_type": "System"
}
```

---

## 6. 接口设计

### 6.1 Transaction Service

Proto Service: `BorrowerPropertyTransaction`

#### Create Transaction

```
POST /api/v1/borrower_property_transaction
RPC: BorrowerPropertyTransaction.Create
```

| 请求字段 | 类型 | 必填 | 说明 |
|----------|------|------|------|
| borrower_property_id | int64 | yes | 关联 Property |
| type | string | yes | Transaction Type |

#### Detail Transaction

```
GET /api/v1/borrower_property_transaction/:id
RPC: BorrowerPropertyTransaction.Detail
```

#### Update Transaction

```
PUT /api/v1/borrower_property_transaction/:id
RPC: BorrowerPropertyTransaction.Update
```

可编辑字段: coe_date, purchase_amount, emd_amount, property_type, pm, strategy, note

#### Delete Transaction

```
DELETE /api/v1/borrower_property_transaction/:id
RPC: BorrowerPropertyTransaction.Delete
```

前置条件: Transaction 未绑定 Loan

#### Switch Status

```
PUT /api/v1/borrower_property_transaction/:id/status
RPC: BorrowerPropertyTransaction.SwitchStatus
```

#### All Transactions

```
GET /api/v1/borrower_property_transaction/all
RPC: BorrowerPropertyTransaction.All
```

#### Update Checklist

```
PUT /api/v1/borrower_property_transaction/:id/checklist
RPC: BorrowerPropertyTransaction.UpdateChecklist
```

手动确认 all_checked，统计字段保留原值。

#### Page Transaction（分页搜索）

```
POST /api/v1/borrower_property_transaction/page
RPC: BorrowerPropertyTransaction.Page
```

| Filter 字段 | 类型 | 必填 | 说明 |
|-------------|------|------|------|
| keyword | string | no | 模糊匹配 Property 地址 (address1/city/state) |
| type | string | no | 精确匹配 Transaction Type |
| status | string | no | 精确匹配 Transaction Status |

响应: id, type, status, coe_date, created_at, updated_at, borrower_property(含完整 address 对象)

#### Operation Stats

```
GET /api/v1/property-transaction/:id/operations/stats
RPC: BorrowerPropertyTransaction.OperationStats
```

| 响应字段 | 类型 | 说明 |
|----------|------|------|
| total | int64 | 操作记录总数 |
| items | []OperationStatsItem | 各类型操作数量 (type + count) |

#### Operation List

```
GET /api/v1/property-transaction/:id/operations
RPC: BorrowerPropertyTransaction.OperationList
```

| 请求字段 | 类型 | 必填 | 说明 |
|----------|------|------|------|
| operation_type | string | no | 按操作类型筛选 |

响应: 操作记录列表，按 created_at DESC 排序

### 6.2 Document Service

Proto Service: `Document`

| 接口 | Method | 路由 | 说明 |
|------|--------|------|------|
| Detail | GET | `/api/v1/document/tree/:entity_type/:entity_id` | 独立路由，返回递归文档树 |
| LeafList | GET | `/api/v1/document/leaf/:entity_type/:entity_id` | 独立路由，叶子节点列表 |
| CreateDocument | POST | `/api/v1/borrower_property_transaction/:tid/document` | 嵌套路由 |
| UploadFile | POST | `/api/v1/borrower_property_transaction/:tid/document/:did/file` | 嵌套路由 |
| DeleteFile | DELETE | `/api/v1/borrower_property_transaction/:tid/document/:did/file/:fid` | 嵌套路由 |

**路由策略**: 读操作使用独立路由（通用，可复用于其他 entity_type），写操作嵌套在 transaction 路由下（便于权限校验和操作日志记录）。

### 6.3 Loan 变更

| 接口 | 变更 | 说明 |
|------|------|------|
| Apply Loan | 新增字段 | borrower_property_transaction_id，XOR 校验 |
| Save Draft | 新增字段 | borrower_property_transaction_id，灵活策略 |
| Update Draft | 新增字段 | borrower_property_transaction_id |

---

## 7. 风险提示与应对

### 7.1 数据一致性

| 风险 | 等级 | 说明 | 应对 |
|------|------|------|------|
| 文档树递归删除 | 中 | 删除 Transaction 需递归清理 Document + File，可能遗漏 | 使用事务包裹，递归删除前验证 Loan 绑定状态 |
| Checklist 统计不准 | 中 | 并发文件操作可能导致统计偏差 | 当前阶段并发低，write-on-update 策略可接受。后续可加乐观锁 |
| Code 碰撞 | 低 | base62(毫秒时间戳) 理论上可能碰撞 | UNIQUE 索引兜底，插入失败重试 |

### 7.2 性能

| 风险 | 等级 | 说明 | 应对 |
|------|------|------|------|
| 文档树递归查询 | 中 | Detail 接口递归加载整棵树，节点多时可能慢 | 模板结构固定（~50节点），深度有限（3层），当前可接受 |
| 分页搜索 JOIN | 低 | LEFT JOIN borrower_properties 匹配地址 | JSONB 地址字段无索引，大数据量时需评估 GIN 索引 |
| Checklist 递归计算 | 低 | 每次文件操作触发全树遍历 | 文件操作频率低，树深度有限 |

### 7.3 兼容性

| 风险 | 等级 | 说明 | 应对 |
|------|------|------|------|
| 现有 Loan 流程 | 中 | Apply/Draft 接口变更影响现有逻辑 | 新字段 DEFAULT 0，向后兼容；XOR 校验仅在有 transaction 时生效 |
| 现有 loan_operations | 低 | 新 operations 表与旧表并存 | 不迁移旧数据，新旧表独立使用 |
| File 表扩展 | 低 | 新增列 DEFAULT 值，不影响现有数据 | ALTER TABLE ADD COLUMN DEFAULT 无锁 |

### 7.4 运维

| 风险 | 等级 | 说明 | 应对 |
|------|------|------|------|
| 种子数据更新 | 低 | 模板结构变更需更新种子 SQL | 开发环境重新 seed，生产环境需单独 migration |
| S3 文件清理 | 中 | 删除 File 记录后需同步删除 S3 对象 | 软删除 + 异步清理策略，避免阻塞请求 |

---

## 8. 实施状态

### 8.1 完成统计

| 状态 | Phase 数 | 任务数 |
|------|---------|-------|
| 已完成 | Phase 1-13 | T001-T038 |
| 待完成 | Phase 14 (验证) | T039-T062 |

### 8.2 关键文件索引

| 层 | 路径 |
|----|------|
| Migration | `deploy/zeebee/sql/20260520151433_borrower_transaction.sql` |
| Seeds | `specs/001-borrower-property-transaction/seeds/001_document_templates.sql` |
| Proto | `app/borrower/rpc/pb/borrower.proto` |
| API Desc | `app/borrower/api/desc/borrower_property_transaction.api` |
| Transaction Model | `domains/zeebee/model/borrower_property_transaction.go` |
| Transaction Repo | `domains/zeebee/repo/borrower_property_transaction.go` |
| Transaction DTO | `domains/zeebee/data/dto/borrower_property_transaction.go` |
| Checklist Logic | `domains/zeebee/logic/checklist_logic.go` |
| Operation Logic | `domains/zeebee/logic/operation_logic.go` |
| Transaction Constants | `domains/zeebee/constants/borrower_property_transaction.go` |
| RPC Logic | `app/borrower/rpc/internal/logic/borrowerpropertytransaction/` |
| API Handler | `app/borrower/api/internal/handler/borrowerpropertytransaction/` |

---

## 9. 研究决策摘要

| ID | 决策 | 备选方案 | 拒绝原因 |
|----|------|----------|----------|
| R1 | 新建独立 operations 表 | 复用 loan_operations | 通用设计，支持多 entity_type |
| R2 | 模板递归创建文档树 | — | — |
| R3 | Presigned URL 直传 S3 | 服务端中转 | 避免服务端瓶颈 |
| R4 | base62(毫秒时间戳)生成 Code | UUID | 更短、可读性更好 |
| R5 | 仅允许叶子节点同级扩展 | 开放任意层级创建 | 保持树深度不变 |
| R6 | 保留旧 loan_operations 不动 | 迁移到 operations | 避免迁移风险 |
| R7 | write-on-update 计算 Checklist | 实时计算 / 定时批处理 | 实时开销大，批处理有延迟 |
| R8 | metadata jsonb 存 is_required | 独立列 | 避免 schema 频繁变更 |
| R9 | common + per-type 模板结构 | 单一模板 | 不同 type 需差异化文档 |
| R10 | Apply XOR / Draft 灵活 | 统一策略 | Apply 需严格，Draft 需灵活 |
| R11 | Stats + List 拆分两个接口 | 单一接口 | 前端两种展示需求 |

---

## 10. 澄清记录

| 日期 | 议题 | 结论 |
|------|------|------|
| 2026-05-19 | 初始需求 | 确定 Transaction 基本功能和 Document 模板结构 |
| 2026-05-20 | Document 树约束 | 允许叶子节点同级追加，checklist 分组统计 |
| 2026-05-22 | 分页搜索设计 | property 地址返回完整对象，keyword 模糊匹配 address1/city/state |
| 2026-05-25 | Loan 二选一 | Apply XOR，Draft 灵活 |
| 2026-05-26 | 操作记录 API | 拆分 Stats + List，模板使用 common + per-type map |
