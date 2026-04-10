'use client';

import { useEffect, useRef, useState } from 'react';
import { Drawer, Form, Input, InputNumber, Switch, Button, message, Space } from 'antd';
import { createRequest } from '@repo/shared-methods/utils';

const { TextArea } = Input;

const DiscoverSettingsDrawer = ({ open, onClose }) => {
	const { request } = createRequest();
	const requestRef = useRef(request);
	useEffect(() => {
		requestRef.current = request;
	});

	const [form] = Form.useForm();
	const [loading, setLoading] = useState(false);
	const [saving, setSaving] = useState(false);

	useEffect(() => {
		if (!open) return;
		let cancelled = false;
		(async () => {
			try {
				setLoading(true);
				const res = await requestRef.current('/lender_discover/settings');
				if (!cancelled) {
					form.setFieldsValue(res?.data || {});
					setLoading(false);
				}
			} catch {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [open, form]);

	const onSave = async () => {
		try {
			const values = await form.validateFields();
			setSaving(true);
			await request('/lender_discover/settings', {
				method: 'PUT',
				body: values,
			});
			message.success('Discover settings saved.');
			setSaving(false);
			onClose?.();
		} catch {
			setSaving(false);
		}
	};

	return (
		<Drawer
			title="Lender discover (OpenAI)"
			width={520}
			open={open}
			onClose={onClose}
			extra={
				<Space>
					<Button onClick={onClose}>Close</Button>
					<Button type="primary" loading={saving} onClick={onSave}>
						Save
					</Button>
				</Space>
			}
		>
			<Form form={form} layout="vertical" disabled={loading}>
				<Form.Item
					name="schedule_enabled"
					label="Enable scheduled search"
					valuePropName="checked"
				>
					<Switch />
				</Form.Item>
				<Form.Item
					name="interval_days"
					label="Interval (days)"
					rules={[{ required: true, message: 'Required' }]}
				>
					<InputNumber min={1} className="w-full" />
				</Form.Item>
				<Form.Item
					name="scope_state"
					label="State (leave empty for nationwide)"
					rules={[{ max: 8, message: 'Max 8 chars' }]}
				>
					<Input maxLength={8} placeholder="e.g. CA (optional)" />
				</Form.Item>
				<Form.Item name="scope_extra" label="Scope (extra description for the model)">
					<TextArea rows={3} />
				</Form.Item>
				<Form.Item name="user_prompt" label="Additional prompt instructions">
					<TextArea rows={4} />
				</Form.Item>
				<Form.Item name="openai_model" label="OpenAI model">
					<Input placeholder="gpt-4o-mini" />
				</Form.Item>
			</Form>
			<p className="text-xs text-base-label mt-2">
				Set <code className="text-xs">OPENAI_API_KEY</code> on the admin server; it is never
				stored in the UI or database.
			</p>
		</Drawer>
	);
};

export default DiscoverSettingsDrawer;
