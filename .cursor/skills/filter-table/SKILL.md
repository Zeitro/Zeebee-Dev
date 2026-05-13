---
name: filter-table
description: >-
  Zeebee 前端「筛选 + 表格」列表页：优先 FilterTable / FilterTableCard（@repo/shared-methods），
  勿手写 FilterField + NewTable。用于 borrower / admin / lender 中带检索的数据列表。
---

# filter-table：筛选 + 表格列表页

`FilterField`、`NewTable` 为原子组件；`FilterTable` / `FilterTableCard` 已封装筛选、分页、请求与 ref（如 `handleSearch`）。标准列表页应优先用封装组件，与 `LoansPage` 等现网一致。

## 何时用哪个

| 场景 | 使用 |
|------|------|
| 筛选 + **表格** + 后端分页（最常见） | **`FilterTable`** |
| 需 **表格/卡片** 切换或卡片列表 | **`FilterTableCard`**（`typeKey`、`renderItem`） |
| 纯静态 demo、且与统一分页协议无关 | 可例外；**默认**仍优先尝试 `getData` / `dataSource` + `FilterTable` |

## Checklist

1. **页面壳**：class 与同应用其它列表页一致（如 borrower 的 `filterWrapClassName`，**抄邻页**）。
2. **`ref`**：`ref={tableFilterRef}`，刷新用 `tableFilterRef.current?.handleSearch?.()`（以组件 `useImperativeHandle` 为准）。
3. **`filterList`**：`filterType`、`paramsKey`、按钮里 `handleSearch`，对齐 `LoansPage` / `management/*/page`。
4. **`api` / `getData`**：有接口传 `api` + `method`；mock 可在 `getData` 返回 `{ list, total }`，并设好 `hasInitRequest` / `dataSource`。
5. **列与刷新**：`columns` 自页面或 `useColumn`；删改后 **`successCallback` → handleSearch**。
6. **列内日期/金额等**：先查 **`@repo/shared-methods/utils`**（**`lib/date.js`**、**`lib/format.js`**），优先 **`formatDate` / `formatTime` / `formatAmount` / `formatPhone`**，勿重复造格式化逻辑。
7. **避免**：标准 CRUD 页只拼 `FilterField` + `NewTable` 并自管分页/请求。

## 参考

- `frontend/borrower/src/app/loans/components/ListView/LoansPage.js` — `FilterTable` + `api` + ref
- `frontend/borrower/src/app/management/borrower/page.js` — `FilterTableCard`
- `frontend/packages/shared-methods/src/components/FilterTable/index.js`
- `frontend/packages/shared-methods/src/components/FilterTableCard/index.js`

## myTransactions 类页面

若将接统一分页接口，应改为 **`FilterTable`（或 `FilterTableCard`）**，避免与全站列表行为分叉。
