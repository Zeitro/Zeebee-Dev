'use client';

import React, { useEffect, useMemo, useState } from 'react';
import { Form, message, Select } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';

/** 与列表页 `TRANSACTION_LIST_API` 同前缀的创建接口；若后端路径不同，只改此处 */
const TRANSACTION_CREATE_API = '/transaction';

const TRANSACTION_TYPE_OPTIONS = [
	{ label: 'Purchase', value: 'Purchase' },
	{ label: 'Sale', value: 'Sale' },
	{ label: 'Re-finance', value: 'Re-finance' },
	{ label: 'Change Title', value: 'Change Title' },
	{ label: 'Assignment', value: 'Assignment' },
];

function formatPropertyOptionLabel(p) {
	const a = p?.address;
	const line1 = a?.address1 ?? '';
	const city = a?.city ?? '';
	const state = a?.state ?? '';
	const tail = [line1, city, state].filter(Boolean).join(', ');
	return tail ? `P-${p?.id} - ${tail}` : `P-${p?.id}`;
}

/**
 * @param {object} props
 * @param {boolean} props.open
 * @param {() => void} props.onClose
 * @param {() => void} [props.onSaved] 保存成功后回调（如刷新表格）
 * @param {(values: { property_id: number, type: string }) => Promise<void>} [props.onSubmit] 自定义提交；不传则 POST `TRANSACTION_CREATE_API`
 * @param {object} [props.initialProperty] 预填物业（如从 Properties 列表行传入）；有值时不再请求 `/property/all`，且物业下拉只读
 */
const NewTransactionModal = ({ open, onClose, onSaved, onSubmit, initialProperty }) => {
	const [form] = Form.useForm();
	const { request } = useMemo(() => createRequest(), []);
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [remotePropertyLoading, setRemotePropertyLoading] = useState(false);
	const [remoteProperties, setRemoteProperties] = useState([]);
	const isPresetProperty = Boolean(initialProperty?.id);
	const propertyLoading = isPresetProperty ? false : remotePropertyLoading;

	useEffect(() => {
		if (!open) return;
		form.resetFields();
		if (isPresetProperty) {
			form.setFieldsValue({ property_id: initialProperty.id });
			return;
		}

		let cancelled = false;
		(async () => {
			await Promise.resolve();
			if (cancelled) return;
			try {
				setRemotePropertyLoading(true);
				const res = await request('/property/all', { method: 'GET' });
				if (!cancelled) setRemoteProperties(res?.lists ?? []);
			} catch {
				if (!cancelled) setRemoteProperties([]);
			} finally {
				if (!cancelled) setRemotePropertyLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [open, form, request, initialProperty, isPresetProperty]);

	const propertySelectOptions = useMemo(() => {
		const source = isPresetProperty ? [initialProperty] : remoteProperties;
		return (source ?? []).map((p) => ({
			value: p.id,
			label: formatPropertyOptionLabel(p),
		}));
	}, [initialProperty, isPresetProperty, remoteProperties]);

	const defaultSubmit = async (values) => {
		await request(TRANSACTION_CREATE_API, {
			method: 'POST',
			body: {
				property_id: values.property_id,
				type: values.type,
			},
		});
	};

	const handleSave = async () => {
		try {
			const values = await form.validateFields();
			setConfirmLoading(true);
			if (onSubmit) {
				await onSubmit(values);
			} else {
				await defaultSubmit(values);
			}
			message.open({
				type: 'success',
				content: 'Transaction created successfully',
			});
			onSaved?.();
			onClose();
		} catch (e) {
			if (e?.errorFields) return;
			console.error(e);
		} finally {
			setConfirmLoading(false);
		}
	};

	if (!open) return null;

	return (
		<NewModal
			width={520}
			title="New Transaction"
			open={open}
			okText="Save"
			cancelText="Cancel"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={onClose}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				disabled={confirmLoading}
				scrollToFirstError
			>
				<Form.Item
					label="Property"
					name="property_id"
					rules={[{ required: true, message: 'Please select a property' }]}
				>
					<Select
						placeholder="Select property"
						allowClear={!initialProperty?.id}
						showSearch={!initialProperty?.id}
						disabled={!!initialProperty?.id}
						loading={propertyLoading}
						options={propertySelectOptions}
						optionFilterProp="label"
					/>
				</Form.Item>
				<Form.Item
					label="Type"
					name="type"
					rules={[{ required: true, message: 'Please select a type' }]}
				>
					<Select
						placeholder="Select type"
						allowClear
						options={TRANSACTION_TYPE_OPTIONS}
					/>
				</Form.Item>
			</Form>
		</NewModal>
	);
};

export default NewTransactionModal;
export { TRANSACTION_TYPE_OPTIONS };
