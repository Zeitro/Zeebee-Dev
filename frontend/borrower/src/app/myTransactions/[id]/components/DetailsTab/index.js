'use client';

/**
 * 交易详情 Details Tab：Property Info / Transaction Info / Title/Escrow 三模块。
 * 编辑态由 DetailsForm 承载；只读态按 constants 字段配置渲染。
 * Property Info 数值字段初始值取自 rowData.borrower_property。
 */

import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Button, Form, message } from 'antd';
import { EditOutlined, SaveOutlined } from '@ant-design/icons';
import { createRequest } from '@repo/shared-methods/utils';
import DetailField from './DetailField';
import DetailsForm from './DetailsForm';
import {
	EMPTY_ESCROW,
	PROPERTY_INFO_VIEW_FIELDS,
	TITLE_ESCROW_VIEW_ROWS,
	TRANSACTION_INFO_VIEW_ROWS,
} from './constants';
import { PermissionSlot } from '@repo/shared-methods';
import { PERMISSIONS } from '@config/permissions';

/**
 * @param {Array} fields
 * @param {object|null} data
 * @param {string} [gridClassName]
 */
function renderViewFields(
	fields,
	data,
	gridClassName = 'grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3',
) {
	return (
		<div className={gridClassName}>
			{fields.map((field) => (
				<DetailField key={field.key} label={field.label} value={field.getValue(data)} />
			))}
		</div>
	);
}

/**
 * @param {Array<{ gridClassName: string, fields: Array }>} rows
 * @param {object|null} data
 */
function renderViewRows(rows, data) {
	return (
		<div className="flex flex-col gap-6">
			{rows.map((row, index) => (
				<div key={index} className={row.gridClassName}>
					{row.fields.map((field) => (
						<DetailField
							key={field.key}
							label={field.label}
							value={field.getValue(data)}
						/>
					))}
				</div>
			))}
		</div>
	);
}

/**
 * Title/Escrow 只读：按 escrows 数组逐条渲染；无数据时展示一行空值。
 * @param {Array|null|undefined} escrows
 */
function renderEscrowViews(escrows) {
	const list = Array.isArray(escrows) && escrows.length > 0 ? escrows : [{}];
	return (
		<div className="flex flex-col gap-6">
			{list.map((escrow, escrowIndex) => (
				<div key={escrow?.id ?? escrowIndex} className="flex flex-col gap-6">
					{TITLE_ESCROW_VIEW_ROWS.map((row, rowIndex) => (
						<div key={rowIndex} className={row.gridClassName}>
							{row.fields.map((field) => (
								<DetailField
									key={field.key}
									label={field.label}
									value={field.getValue(escrow)}
								/>
							))}
						</div>
					))}
				</div>
			))}
		</div>
	);
}

/**
 * @param {{ title: string, extra?: React.ReactNode, bordered?: boolean, children: React.ReactNode }} props
 */
function DetailModule({ title, extra, bordered = true, children }) {
	return (
		<section
			className={['flex flex-col gap-4', bordered ? 'border-b border-base-border pb-8' : '']
				.filter(Boolean)
				.join(' ')}
		>
			<div className="flex items-center justify-between gap-4">
				<h2 className="m-0 text-xl font-medium leading-7 text-foreground">{title}</h2>
				{extra}
			</div>
			{children}
		</section>
	);
}

const DetailsTab = ({ rowData, transactionId, onRefresh, isSearchParamsEditing }) => {
	const [form] = Form.useForm();
	const [messageApi, contextHolder] = message.useMessage();
	const { request } = useMemo(() => createRequest(), []);
	const [isEditing, setIsEditing] = useState(isSearchParamsEditing);
	const [saving, setSaving] = useState(false);

	const syncFormValues = useCallback(() => {
		const escrows = rowData?.escrows?.length > 0 ? rowData.escrows : [EMPTY_ESCROW];
		form.setFieldsValue({ ...rowData, escrows });
	}, [form, rowData]);

	useEffect(() => {
		onRefresh();
	}, [onRefresh]);

	useEffect(() => {
		if (isEditing) {
			syncFormValues();
		}
	}, [isEditing, syncFormValues]);

	const handleEdit = () => {
		syncFormValues();
		setIsEditing(true);
	};

	const handleSave = async () => {
		try {
			setSaving(true);
			const values = await form.validateFields();
			values.property.id = rowData?.property?.id;
			await request(`/borrower_property_transaction/${transactionId}`, {
				method: 'PUT',
				body: {
					...values,
				},
			});
			setIsEditing(false);
			messageApi.success('Transaction updated');
			onRefresh?.();
		} catch {
			// validation or request errors are handled by antd / request layer
		} finally {
			setSaving(false);
		}
	};

	const headerAction = isEditing ? (
		<Button type="primary" icon={<SaveOutlined />} loading={saving} onClick={handleSave}>
			Save
		</Button>
	) : (
		<Button icon={<EditOutlined />} onClick={handleEdit}>
			Edit
		</Button>
	);

	const editAction = (
		<PermissionSlot code={PERMISSIONS?.TRANSACTION?.EDIT}>{headerAction}</PermissionSlot>
	);

	return (
		<div className="flex flex-col gap-8 rounded-lg border border-base-border bg-background px-8 py-6">
			{contextHolder}
			{isEditing ? (
				<Form form={form} layout="vertical" scrollToFirstError size="large">
					<DetailsForm rowData={rowData} propertyInfoExtra={editAction} saving={saving} />
				</Form>
			) : (
				<>
					<DetailModule title="Property Info" extra={editAction}>
						{renderViewFields(PROPERTY_INFO_VIEW_FIELDS, rowData)}
					</DetailModule>

					<DetailModule title="Transaction Info">
						{renderViewRows(TRANSACTION_INFO_VIEW_ROWS, rowData)}
					</DetailModule>

					<DetailModule title="Title/Escrow" bordered={false}>
						{renderEscrowViews(rowData?.escrows)}
					</DetailModule>
				</>
			)}
		</div>
	);
};

export default DetailsTab;
