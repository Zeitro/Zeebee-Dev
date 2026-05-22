import React from 'react';
import { UploadModalBtn } from '@repo/shared-methods/componentsUi/zeeBeeDetails/components/AddDocument/UploadModal';
import { displayValue, formatDate } from '@repo/shared-methods/utils';
import { FolderOpenOutlined } from '@ant-design/icons';

/**
 * CollapsibleDocumentList 列配置（交易 Documents Tab）
 *
 * rowPlacement:
 * - parent：仅父级（文档类型行）
 * - file：仅子级（文件行）；父级该列空占位对齐
 * - both：父子都展示；可用 parentDataIndex / fileDataIndex 区分字段
 *
 * 渲染：
 * - render(value, record, { level, dataIndex })：仅父级（both 列的父级侧也用 render，子级勿用）
 * - fileRender(value, defaultNode, file, parentRecord)：子级或 both 的子级侧
 *   - 省略 fileRender 时：首列默认文件名 pill，末列默认 View / Delete
 * - formatter：无自定义 render 时用于 displayValue(value, formatter)
 */
export const transactionDocumentColumns = () => [
	{
		title: 'Document type',
		key: 'document_type',
		dataIndex: 'document_type',
		minWidth: 200,
		ellipsis: true,
		rowPlacement: 'both',
		// 父级展示类型名；子级走组件默认（文件名 + 图标 pill），无需 fileRender
	},
	{
		title: 'Update date',
		key: 'updated_at',
		width: 120,
		ellipsis: true,
		rowPlacement: 'file',
		fileDataIndex: 'updated_at',
		fileRender: (value, _defaultNode, file, parentRecord) =>
			displayValue(value ?? file?.updated_at ?? parentRecord?.updated_at, formatDate),
	},
	{
		title: 'Status',
		key: 'status',
		width: 150,
		ellipsis: true,
		rowPlacement: 'file',
		fileDataIndex: 'status',
		fileRender: (value, _defaultNode, file, parentRecord) =>
			displayValue(value ?? file?.status ?? parentRecord?.status),
	},
	{
		title: '',
		key: 'operation',
		width: 120,
		align: 'right',
		rowPlacement: 'both',
		render: (_, record) => (
			<UploadModalBtn isChangeType={false} record={record} isStopPropagation>
				Upload
			</UploadModalBtn>
		),
		// 子级省略 fileRender → 默认 View / Delete（DefaultFileActions）
	},
];
