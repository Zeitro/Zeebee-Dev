'use client';

import React, { useCallback, useEffect, useState } from 'react';
import { Button, Form } from 'antd';
import {
	EditOutlined,
	EnvironmentOutlined,
	HomeOutlined,
	AccountBookOutlined,
	SaveOutlined,
} from '@ant-design/icons';
import { createRequest } from '@repo/shared-methods/utils';
import DetailField from './DetailField';
import DetailSection from './DetailSection';
import PropertyDetailsForm from './PropertyDetailsForm';
import { ADDRESS_FIELDS, LOAN_INFO_FIELDS, PROPERTY_FIELDS } from './constants';

function renderFields(fields, data) {
	return fields.map((field) => (
		<DetailField key={field.label} label={field.label} value={field.getValue(data)} />
	));
}

const PROPERTY_FORM_FIELDS = ['property_type', 'living_area_sqft', 'unit', 'bedrooms', 'bathrooms'];

const DetailsTab = ({ rowData, propertyId, canEdit, onRefresh }) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [isEditing, setIsEditing] = useState(false);
	const [saving, setSaving] = useState(false);

	const syncFormValues = useCallback(() => {
		form.setFieldsValue({
			property_type: rowData?.property_type,
			living_area_sqft: rowData?.living_area_sqft,
			unit: rowData?.unit,
			bedrooms: rowData?.bedrooms,
			bathrooms: rowData?.bathrooms,
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
			const values = await form.validateFields(PROPERTY_FORM_FIELDS);
			await request(`/property/${propertyId}`, {
				method: 'PUT',
				body: {
					...rowData,
					...values,
				},
			});
			setIsEditing(false);
			onRefresh?.();
		} catch {
			// validation or request errors are handled by antd / request layer
		} finally {
			setSaving(false);
		}
	};

	const headerAction =
		canEdit &&
		(isEditing ? (
			<Button type="primary" icon={<SaveOutlined />} loading={saving} onClick={handleSave}>
				Save
			</Button>
		) : (
			<Button icon={<EditOutlined />} onClick={handleEdit}>
				Edit
			</Button>
		));

	return (
		<div className="rounded-lg border border-base-border bg-background p-4">
			<DetailSection icon={<EnvironmentOutlined />} title="Address" extra={headerAction}>
				<div className="grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
					{renderFields(ADDRESS_FIELDS, rowData)}
				</div>
			</DetailSection>

			<DetailSection icon={<HomeOutlined />} title="Property details">
				{isEditing ? (
					<Form form={form} layout="vertical" disabled={saving} scrollToFirstError>
						<PropertyDetailsForm rowData={rowData} />
					</Form>
				) : (
					<div className="grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
						{renderFields(PROPERTY_FIELDS, rowData)}
					</div>
				)}
			</DetailSection>

			<DetailSection icon={<AccountBookOutlined />} title="Loan information" bordered={false}>
				<div className="grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
					{renderFields(LOAN_INFO_FIELDS, rowData)}
				</div>
			</DetailSection>
		</div>
	);
};

export default DetailsTab;
