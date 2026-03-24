'use client';
import { NewModal } from '@repo/shared-methods';
import { Form } from 'antd';
import { useState } from 'react';
import dynamic from 'next/dynamic';
// 动态导入避免 SSR 问题
const MonacoEditor = dynamic(() => import('@repo/shared-methods/components/MonacoEditor'), {
	ssr: false,
});
export const MonacoModal = ({ setOpen, data }) => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [value, setValue] = useState(data || '');
	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const values = await form.validateFields();
			setConfirmLoading(false);
			setOpen(false);
		} catch (errpr) {
			setConfirmLoading(false);
		}
	};

	return (
		<NewModal
			width={900}
			title="Add variable"
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<MonacoEditor language={'html'} value={value} onChange={(value) => setValue(value)} />
		</NewModal>
	);
};
