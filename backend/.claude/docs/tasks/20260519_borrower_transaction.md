# Borrower Property Transaction 需求

## 需求描述

为 borrower property 添加 transaction 管理，每个 transaction 关联通用文档管理系统。在申请 loan 前准备 transaction
相关文件，作为资产管理的一部分。

## 开发分支

`feat/backend/borrower/property-transaction`

## Scope

`backend/borrower`

---

## 功能需求

### 1. Transaction 管理

为 property 新增 transaction，一个 property 可创建多个 transaction（允许 type 相同）。

| 操作            | 描述                                                                          |
|---------------|-----------------------------------------------------------------------------|
| Create        | 选择关联 property + type，自动生成 code，自动创建对应 documents                             |
| Detail        | 根据 id 获取详情                                                                  |
| Edit          | 编辑 coe_data, purchase_amount, emd_amount, property_type, pm, strategy, note |
| Delete        | 只删除未绑定 loan 的 transaction，递归删除关联 document 及 file                            |
| Switch Status | Pending / Closed / Cancelled 三状态任意切换                                        |
| All           | 用于 select 下拉的 options                                                       |
| Checklist     | 更新 checklist 字段确认                                                           |

**Transaction Type**: `Purchase`, `Sale`, `Re-finance`, `Change Title`, `Assignment`

**Transaction Status**: `Pending`, `Closed`, `Cancelled`

**Code 生成规则**: `tc` + 年份(2位) + base62 加密的毫秒时间戳，例: `tc26ttdhd10`

**可编辑字段约束**:

- `property_type` 与 `borrower_properties.property_type` 类型一致
- `strategy` 与 `loans.business_strategy` 字段一致（单值）

**Checklist 结构**:

```json
{
  "pending": { "all_checked": true, "document_id": 1 },
  "closed": { "all_checked": false, "document_id": 0 },
  "cancelled": { "all_checked": false, "document_id": 0 }
}
```

### 2. Document 管理（通用文档系统）

新增通用文档管理，document 为文件夹树结构，file 通过 `entity_type` + `entity_id` 多态关联。

**架构关系**:

```
business (transaction) → entity_documents (绑定) → documents (文件夹树) → files (文件实体)
```

| 操作              | 描述                                                                                              |
|-----------------|-------------------------------------------------------------------------------------------------|
| Create Document | 只创建最底层同级 document，同 level 内 name 不重复                                                            |
| Upload File     | 只上传到最底层 document，storage_type=Aws S3, upload_type=Manual，更新 document 及 parent 的 last_updated_at |
| Delete File     | 只删除 upload_type=Manual 的 file                                                                   |

**Document 模板**: 创建 transaction 时从 `document_templates` 按 type 加载 structure，递归创建文档目录树。模板可复用于其他业务。
`document_templates` 不做 CRUD 管理，数据手动配置入库。模板结构见附录 A。

### 3. Loan 关联

创建 loan 时可选绑定 transaction，通过 transaction 间接获取 property。property 与 transaction 二选一。

### 4. 操作日志

新增 `operations` 表，统一记录 transaction 和 loan 的操作日志。

---

## 数据模型

### borrower_property_transactions

| 字段                               | 类型                 | 说明                                               |
|----------------------------------|--------------------|--------------------------------------------------|
| id                               | serial PK          |                                                  |
| borrower_property_id             | int4               | 关联 property                                      |
| organization_id                  | int4               | 组织 ID                                            |
| code                             | varchar(20) UNIQUE | 自动生成的唯一标识                                        |
| type                             | varchar(20)        | Purchase/Sale/Re-finance/Change Title/Assignment |
| status                           | varchar(20)        | Pending/Closed/Cancelled                         |
| coe_data                         | date               | close of end date                                |
| purchase_amount                  | decimal(10,2)      |                                                  |
| emd_amount                       | decimal(10,2)      |                                                  |
| property_type                    | varchar(50)        | 与 borrower_properties.property_type 一致           |
| pm                               | varchar(255)       |                                                  |
| strategy                         | varchar(100)       | 与 loans.business_strategy 一致                     |
| note                             | varchar(255)       |                                                  |
| checklist                        | jsonb              | 确认状态                                             |
| created_at / updated_at / is_del |                    | 标准字段                                             |

索引: `idx_bp_trans_property_id`, `idx_bp_trans_org_id`, `idx_bp_trans_code(UNIQUE)`

### operations（操作日志）

| 字段                               | 类型           | 说明      |
|----------------------------------|--------------|---------|
| id                               | serial PK    |         |
| entity_id                        | int4         | 业务实体 ID |
| entity_type                      | varchar(100) | 业务实体类型  |
| operator_id                      | int4         | 操作人 ID  |
| operator_name                    | varchar(50)  |         |
| operator_type                    | varchar(50)  |         |
| operation                        | varchar(50)  | 操作类型    |
| note                             | varchar(255) |         |
| metadata                         | jsonb        |         |
| organization_id                  | int4         |         |
| created_at / updated_at / is_del |              |         |

索引: `idx_operation_entity(entity_id, entity_type)`, `idx_operation_org_id`

### entity_documents（业务-文档绑定）

| 字段                               | 类型          | 说明     |
|----------------------------------|-------------|--------|
| id                               | serial PK   |        |
| entity_type                      | varchar(20) | 业务类型   |
| entity_id                        | int4        | 业务 ID  |
| root_document_id                 | int4        | 根文档 ID |
| created_at / updated_at / is_del |             |        |

索引: `idx_entity_doc_entity_id(entity_type, entity_id)`, `idx_entity_doc_root_id`

### documents（文档树）

| 字段                                | 类型           | 说明  |
|-----------------------------------|--------------|-----|
| id                                | serial PK    |     |
| parent_id                         | int4         | 父节点 |
| type                              | varchar(20)  |     |
| name                              | varchar(255) |     |
| level                             | int4         | 层级  |
| sort_order                        | int4         | 排序  |
| last_uploaded_at                  | timestamptz  |     |
| metadata                          | jsonb        |     |
| created_at / uploaded_at / is_del |              |     |

索引: `idx_doc_pid`, `idx_doc_name`

### document_templates（文档模板）

| 字段                               | 类型           | 说明    |
|----------------------------------|--------------|-------|
| id                               | serial PK    |       |
| name                             | varchar(255) |       |
| type                             | varchar(20)  |       |
| structure                        | json         | 文档树结构 |
| status                           | varchar(20)  |       |
| created_at / updated_at / is_del |              |       |

### files 表变更

新增列:

- `file_name` varchar(255) — 文件名
- `storage_type` varchar(20) DEFAULT 'Aws S3' — 存储类型 (Aws S3, Google Drive)
- `external_id` varchar(255) — 外部 ID
- `upload_type` varchar(20) DEFAULT 'Manual' — 上传方式 (Manual, Google Drive)

索引: `idx_files_external_id`

### loans 表变更

新增列:

- `borrower_property_transaction_id` int4 DEFAULT 0 — 关联 transaction

索引: `idx_loan_borrower_property_trans_id`

---

## 涉及层次

| 层次        | 变更                                                                                                  |
|-----------|-----------------------------------------------------------------------------------------------------|
| Model     | 新增 BorrowerPropertyTransaction, Document, EntityDocument, Operation, DocumentTemplate；修改 File, Loan |
| Repo      | 新增 TransactionRepo, DocumentRepo, EntityDocumentRepo, OperationRepo, DocumentTemplateRepo           |
| Constants | 新增 TransactionType/Status, DocumentType, FileEntityType, StorageType, UploadType                    |
| Proto     | 新增 Transaction service + messages, Document service + messages                                      |
| API Desc  | 新增 transaction.api, document.api                                                                    |
| RPC Logic | 新增 transaction/, document/ 逻辑                                                                       |
| API Logic | 新增 borrower_transaction/, document/ 逻辑                                                              |

---

## 附录 A: Document 模板结构

```
Pending (sort_order: 1)
├── Contract (sort_order: 1)
│   ├── Original RPA
│   ├── Assignment Agreement AOAA
│   ├── Counter Offer
│   ├── Addendum (Extension)
│   ├── Seller financing Docs
│   └── JV Agreement
├── Escrow (sort_order: 2)
│   ├── Wire Instructions
│   ├── EMD Payment Proof
│   ├── EMD Receipt
│   ├── Prelim Report
│   ├── NHD Report
│   ├── Colored map
│   ├── FinCen
│   ├── Buyer Open Package
│   └── LOI
├── Disclosures (sort_order: 3)
├── Loan (sort_order: 4)
│   ├── Application Form
│   ├── Loan Document
│   ├── Contact form
│   └── Appraisal Report
├── Insurance (sort_order: 5)
│   ├── EOI
│   ├── Payment Receipt
│   ├── Full Policy
│   └── Flood Insurance (Optional)
├── Closing Package (sort_order: 6)
│   ├── Buyer Estimated Statement
│   ├── Final Buyer Statement
│   ├── All Signed Documents
│   ├── Recorded GD
│   └── Recorded DOT
├── 1031 - Optional (sort_order: 7)
│   └── 1031 Exchange Documents
├── Lease - Optional (sort_order: 8)
│   ├── Lease Agreement
│   └── TENANT ESTOPPEL CERTIFICATE
├── HOA - Optional (sort_order: 9)
│   └── See "HOA" Sheet
└── Solar - Optional (sort_order: 10)
    └── See "Solar" Sheet
Closed (sort_order: 2)
Cancelled (sort_order: 3)
```
