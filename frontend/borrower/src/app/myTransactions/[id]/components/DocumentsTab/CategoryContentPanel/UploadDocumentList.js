'use client';

import React, { useMemo, useState } from 'react';
import CollapsibleDocumentList from '@/components/CollapsibleDocumentList';
import { NoData } from '@repo/shared-methods';
import styles from './index.module.scss';
import { createRequest, displayValue, formatDate } from '@repo/shared-methods/utils';
import UploadModal from './uploadModal';
import { Button, Tag, Tooltip, message } from 'antd';
import { GoogleOutlined } from '@ant-design/icons';
import { PlusOutlined, UploadOutlined } from '@ant-design/icons';
import AddRequest from '@repo/shared-methods/componentsUi/zeeBeeDetails/components/AddDocument/AddRequest';

function SectionLabel({ title, count }) {
	return (
		<div className={styles.sectionLabel}>
			<p className={styles.sectionLabelTitle}>{title}</p>
			<span className={styles.sectionLabelCount}>· {count}</span>
		</div>
	);
}

const UploadDocumentList = ({ category, onRefresh, transactionId }) => {
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
	// 列配置不依赖任何响应式值（仅用到稳定的 setState 与导入工具），
	// 直接声明为普通常量，由 React Compiler 自动记忆化。
	const columns = [
		{
			title: 'Document type',
			key: 'name',
			dataIndex: 'name',
			minWidth: 180,
			ellipsis: true,
			rowPlacement: 'both',
			// 父级展示类型名；子级走组件默认（文件名 + 图标 pill），无需 fileRender
		},
		{
			title: 'Update date',
			key: 'updated_at',
			width: 100,
			ellipsis: true,
			rowPlacement: 'file',
			fileDataIndex: 'updated_at',
			fileRender: (value, _defaultNode, file, parentRecord) =>
				displayValue(value ?? file?.updated_at ?? parentRecord?.updated_at, formatDate),
		},
		{
			// 文件已同步到 Google Drive 时展示标签；有标签即表示该文件当前在 Drive 对应文件夹中。
			// 同步日期通过 hover 标签的 tooltip 展示，标签本身不跳转，保持列宽紧凑。
			title: '',
			key: 'drive',
			width: 140,
			rowPlacement: 'file',
			fileDataIndex: 'drive_synced_at',
			fileRender: (value, _defaultNode, file) => {
				const syncedAt = value ?? file?.drive_synced_at;
				if (!syncedAt) {
					return displayValue(null);
				}
				return (
					<Tooltip title={`Synced ${syncedAt}`}>
						<Tag icon={<GoogleOutlined />} color="green" className="m-0">
							Google Drive
						</Tag>
					</Tooltip>
				);
			},
		},
		{
			title: '',
			key: 'operation',
			width: 120,
			align: 'right',
			rowPlacement: 'both',
			render: (_, record) => (
				<Button
					type="default"
					size="small"
					onClick={(e) => {
						e.stopPropagation();
						setRecord(record);
						setOpen(true);
					}}
					icon={<UploadOutlined />}
				>
					Upload
				</Button>
			),
			// 子级省略 fileRender → 默认 View / Delete（DefaultFileActions）
		},
	];
	const [open, setOpen] = useState(false);
	const [record, setRecord] = useState({});
	const [openAdd, setOpenAdd] = useState(false);
	const { request } = useMemo(() => createRequest(), []);
	const [messageApi, contextHolder] = message.useMessage();
	const { requiredList, optionalList } = useMemo(() => {
		const documents = category?.children ?? [];
		return {
			requiredList: documents.filter((d) => d.is_required),
			optionalList: documents.filter((d) => !d.is_required),
		};
	}, [category?.children]);

	const hasRequired = requiredList.length > 0;
	const hasOptional = optionalList.length > 0;

	const deleteFile = async (file) => {
		await request(`/borrower_property_transaction/${transactionId}/file/${file?.id}`, {
			method: 'DELETE',
		});
		messageApi.success('File deleted successfully');
		onRefresh?.();
	};

	return (
		<div className="px-6 py-5">
			{contextHolder}
			<div className="pb-4">
				<Button
					type="primary"
					onClick={() => setOpenAdd(true)}
					size="large"
					icon={<PlusOutlined />}
				>
					Add Document type
				</Button>
			</div>
			{hasRequired && (
				<section className="mt-4">
					<SectionLabel title="Required" count={requiredList.length} />
					<CollapsibleDocumentList
						columns={columns}
						dataSource={requiredList}
						onDeleteFile={deleteFile}
						rowKey="document_id"
					/>
				</section>
			)}

			{hasOptional && (
				<section className={hasRequired ? 'mt-8' : 'mt-4'}>
					<SectionLabel title="Optional" count={optionalList.length} />
					<CollapsibleDocumentList
						columns={columns}
						dataSource={optionalList}
						onDeleteFile={deleteFile}
						rowKey="document_id"
					/>
				</section>
			)}

			{!hasRequired && !hasOptional && (
				<div className="px-6 py-6">
					<NoData />
				</div>
			)}
			{openAdd && (
				<AddRequest
					open={openAdd}
					setOpen={setOpenAdd}
					id={transactionId}
					category={category?.document_id}
					apiType="transaction"
					title="Add document type"
					onSuccessCallback={onRefresh}
				/>
			)}
			{open && (
				<UploadModal
					open={open}
					setOpen={setOpen}
					record={record}
					transaction_id={transactionId}
					onRefresh={onRefresh}
				/>
			)}
		</div>
	);
};

export default UploadDocumentList;
