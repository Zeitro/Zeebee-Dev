'use client';

import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Button, Form, message } from 'antd';
import { EditOutlined, SaveOutlined } from '@ant-design/icons';
import { createRequest } from '@repo/shared-methods/utils';
import DetailField from './DetailField';
import DetailsForm from './DetailsForm';
import { BASIC_INFO_FORM_FIELDS, BASIC_INFO_VIEW_FIELDS } from './constants';

function renderViewFields(fields, data) {
	const gridFields = fields.filter((f) => !f.fullWidth);
	const fullWidthFields = fields.filter((f) => f.fullWidth);

	return (
		<div className="flex flex-col gap-6 rounded-lg p-4">
			<div className="grid w-full grid-cols-1 gap-6 sm:grid-cols-3">
				{gridFields.map((field) => (
					<DetailField key={field.key} label={field.label} value={field.getValue(data)} />
				))}
			</div>
			{fullWidthFields.map((field) => (
				<DetailField
					key={field.key}
					label={field.label}
					value={field.getValue(data)}
					className="w-full"
				/>
			))}
		</div>
	);
}

const DetailsTab = ({ rowData, transactionId, onRefresh }) => {
	const [form] = Form.useForm();
	const { request } = useMemo(() => createRequest(), []);
	const [isEditing, setIsEditing] = useState(false);
	const [saving, setSaving] = useState(false);

	const syncFormValues = useCallback(() => {
		form.setFieldsValue({
			coe_date: rowData?.coe_date ?? rowData?.coeDate,
			purchase_price: rowData?.purchase_price ?? rowData?.purchasePrice,
			emd: rowData?.emd,
			property_type:
				rowData?.property_type ?? rowData?.propertyType ?? rowData?.property?.property_type,
			pm: rowData?.pm,
			business_strategy: rowData?.business_strategy ?? rowData?.businessStrategy,
			tc_notes: rowData?.tc_notes ?? rowData?.tcNotes,
		});
	}, [form, rowData]);

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
			const values = await form.validateFields(BASIC_INFO_FORM_FIELDS);
			await request(`/transaction/${transactionId}`, {
				method: 'PUT',
				body: {
					...rowData,
					...values,
				},
			});
			setIsEditing(false);
			message.success('Transaction updated');
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

	return (
		<div className="flex flex-col gap-6 rounded-lg border border-base-border bg-background px-8 py-6">
			<div className="flex items-center justify-between gap-4">
				<h2 className="m-0 text-xl font-medium leading-7 text-foreground">
					Basic Information
				</h2>
				{headerAction}
			</div>

			{isEditing ? (
				<Form
					form={form}
					layout="vertical"
					disabled={saving}
					scrollToFirstError
					size="large"
				>
					<DetailsForm />
				</Form>
			) : (
				renderViewFields(BASIC_INFO_VIEW_FIELDS, rowData)
			)}
		</div>
	);
};

export default DetailsTab;
