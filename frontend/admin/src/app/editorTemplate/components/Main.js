'use client';
import { Button, Col, Drawer, Form, Input, message, Row, Select, Spin } from 'antd';
import { useCallback, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { VariableBtn } from './Variable';
import { PreviewBtn } from './Preview';
import dynamic from 'next/dynamic';
const { TextArea } = Input;

const MonacoEditor = dynamic(
	() => import('@repo/shared-methods/components/MonacoEditor'), // Adjust the path to your component
	{
		ssr: false, // This is the crucial part
		// Optional: Add a loading component to improve UX
		loading: () => (
			<p className="text-center w-full h-[200px] text-base">Asynchronous loading...</p>
		),
	},
);
//
const MainDrawer = ({ rowData, open, setOpen, successCallback }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	// const variableGroup = Form.useWatch('variable_group_schema', form);
	const templateContent = Form.useWatch('content', form);
	// const variableFields = useMemo(() => {
	// 	let list = [];
	// 	variableGroup?.map?.((item) => {
	// 		list = [...list, ...(item.variables || [])];
	// 	});
	// 	return list;
	// }, [variableGroup]);
	const categoryOptions = [
		{
			label: 'Contract',
			value: 'Contract',
		},
		{
			label: 'Email',
			value: 'Email',
		},
	];
	const systemOptions = [
		{
			label: 'Admin',
			value: 'Admin',
		},
		{
			label: 'Lender',
			value: 'Lender',
		},
		{
			label: 'Borrower',
			value: 'Borrower',
		},
	];
	const scopeOptions = [
		{
			label: 'Organization',
			value: 'Organization',
		},
		{
			label: 'System',
			value: 'System',
		},
	];
	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			let api = '/base_template';
			let method = 'POST';
			if (rowData?.id) {
				api = `/base_template/${rowData?.id}`;
				method = 'PUT';
			}
			const values = await form.validateFields();
			await request(api, { method, body: values });
			message.destroy();
			message.open({
				content: `Template ${rowData?.id ? 'Edited' : 'Added'} successfully.`,
				type: 'success',
			});
			setConfirmLoading(false);
			successCallback?.();
			form.resetFields();
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	const getDetails = useCallback(async () => {
		try {
			setLoading(true);
			const res = await request(`/base_template/${rowData?.id}`);
			form.setFieldsValue(res);
			setLoading(false);
		} catch {
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [rowData?.id]);

	useDeferredSyncEffect({
		predicate: () => rowData?.id,
		deferred: () => {
			getDetails();
		},
		deps: [rowData, getDetails],
		delay: 0,
	});

	return (
		<Drawer
			title={`${rowData?.id ? 'Edit' : 'Add'} editor template`}
			placement="right"
			size="large"
			onClose={() => {
				setOpen(false);
				form.resetFields();
			}}
			open={open}
			extra={
				<Button type="primary" onClick={handleSave} loading={confirmLoading}>
					Save
				</Button>
			}
		>
			<Spin spinning={loading}>
				<Form
					name="editorTemplateForm"
					size="large"
					layout="vertical"
					autoComplete="off"
					form={form}
					className="h-full"
				>
					<Row gutter={16}>
						<Col md={24} xs={24}>
							<Form.Item label="Purpose" name="purpose" rules={[{ required: true }]}>
								<Input placeholder="Purpose" allowClear />
							</Form.Item>
						</Col>
						<Col md={12} xs={24}>
							<Form.Item
								label="Template name"
								name="name"
								rules={[{ required: true }]}
							>
								<Input placeholder="Template name" allowClear />
							</Form.Item>
						</Col>
						<Col md={12} xs={24}>
							<Form.Item
								label="Template category"
								name="category"
								rules={[{ required: true }]}
							>
								<Select
									placeholder="Template category"
									allowClear
									options={categoryOptions}
								/>
							</Form.Item>
						</Col>
						<Col md={12} xs={24}>
							<Form.Item label="System" name="type" rules={[{ required: true }]}>
								<Select placeholder="System" allowClear options={systemOptions} />
							</Form.Item>
						</Col>
						<Col md={12} xs={24}>
							<Form.Item label="Scope" name="scope" rules={[{ required: true }]}>
								<Select placeholder="Scope" allowClear options={scopeOptions} />
							</Form.Item>
						</Col>

						<Col md={12} xs={24}>
							<Form.Item label="Variable" name="variable_group_schema">
								<VariableBtn />
							</Form.Item>
						</Col>
						<Col md={24} xs={24}>
							<Form.Item label="Description" name="description">
								<TextArea
									placeholder="Description"
									allowClear
									autoSize={{ minRows: 2, maxRows: 2 }}
									count={{
										show: true,
									}}
								/>
							</Form.Item>
						</Col>
						<Col md={24} xs={24}>
							<Form.Item label="Example data" name="example_variables">
								<TextArea
									placeholder="Example data"
									allowClear
									autoSize={{ minRows: 2, maxRows: 2 }}
									count={{
										show: true,
									}}
								/>
							</Form.Item>
						</Col>
						<Col xs={24}>
							<PreviewBtn
								className="!absolute right-0 top-0 z-10"
								value={templateContent}
							/>
							<Form.Item label="Template" name="content" rules={[{ required: true }]}>
								{/* <WangEditor
								wrapperProps={{
									className: '!h-[65vh]',
								}}
								variables={variableFields}
								isShowVariablesToolbar={true}
								isShowViewCodeToolbar={true}
							/> */}
								<MonacoEditor height={500} language={'html'} />
							</Form.Item>
						</Col>
					</Row>
				</Form>
			</Spin>
		</Drawer>
	);
};
export default MainDrawer;
