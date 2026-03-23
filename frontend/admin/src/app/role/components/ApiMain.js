import React, { useEffect, useState } from 'react';
import { Row, Col, Form, Select, Input, message, Radio, TreeSelect, InputNumber } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { yesNoOptions, systemTypeOptions, titleCol } from '@repo/shared-methods/enum';
import { createRequest, eventBus } from '@repo/shared-methods/utils';
import useData from './useData';
const { TextArea } = Input;
const ApiMainModal = ({ setOpen, rowData, eventBusPrefix }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const { getApiTree, apiTreeData } = useData();

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/api';
			let method = 'POST';
			let msg = 'Api creation successful.';
			if (rowData?.id) {
				// value.id = rowData?.id;
				url = `/api/${rowData?.id}`;
				method = 'PUT';
				msg = 'Api update successful.';
			}
			await request(url, {
				method,
				body: { ...value },
			});
			eventBus.emit(`${eventBusPrefix}${value?.type}`, rowData?.id ? 'Edit' : 'Add');
			message.open({ type: 'success', content: msg });
			setConfirmLoading(false);
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	// const getList = async (v) => {
	// 	try {
	// 		setLoading(true);
	// 		const payload = new URLSearchParams({ type: v || rowData?.type });
	// 		const api = `/api/tree?${payload}`;
	// 		const res = await request(api);
	// 		setTreeData(res?.list);
	// 		setLoading(false);
	// 	} catch (error) {
	// 		setLoading(false);
	// 		throw error;
	// 	}
	// };

	useEffect(() => {
		if (rowData?.type) {
			const getFn = () => {
				getApiTree(rowData?.type);
			};
			getFn();
		}
	}, [rowData?.type]);

	return (
		<NewModal
			width={600}
			title={`${rowData?.id ? 'Edit' : 'Add'} Api`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				form={form}
				name="api-form"
				layout="vertical"
				autoComplete="off"
				className={[
					' overflow-hidden overflow-y-auto',
					rowData?.id ? 'h-[64vh]' : 'h-[56vh]',
				].join(' ')}
				size="large"
				disabled={confirmLoading}
				scrollToFirstError
				initialValues={{
					status: 'Active',
					auth_required: true,
					permission_required: true,
					sort_order: 6,
					...(rowData || {}),
				}}
			>
				<Row gutter={16}>
					<Col md={18}>
						<Form.Item
							label="Name"
							name={['name']}
							rules={[
								{
									required: true,
								},
							]}
						>
							<Input placeholder="Enter name" />
						</Form.Item>
					</Col>
					<Col md={6}>
						<Form.Item label="Method" name={['method']}>
							<Select
								options={[
									{ value: 'GET', label: 'GET' },
									{ value: 'POST', label: 'POST' },
									{ value: 'PUT', label: 'PUT' },
									{ value: 'DELETE', label: 'DELETE' },
									{ value: 'PATCH', label: 'PATCH' },
								]}
								placeholder="Select"
							/>
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item label="Path" name={['path']}>
							<Input placeholder="Enter path" />
						</Form.Item>
					</Col>
					<Col md={12}>
						<Form.Item
							label="System type"
							name={['type']}
							rules={[
								{
									required: true,
								},
							]}
						>
							<Select
								options={systemTypeOptions}
								placeholder="Select system type"
								onChange={(v) => {
									form.setFieldValue(['parent_id'], null);
									getApiTree(v);
								}}
							/>
						</Form.Item>
					</Col>
					<Col md={12}>
						<Form.Item
							label="Status"
							name={['status']}
							rules={[
								{
									required: true,
								},
							]}
						>
							<Select
								options={[
									{ value: 'Active', label: 'Active' },
									{ value: 'Disable', label: 'Disable' },
								]}
								placeholder="Select status"
							/>
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item label="Parent api" name={['parent_id']}>
							<TreeSelect
								showSearch
								style={{ width: '100%' }}
								styles={{
									popup: { root: { maxHeight: 400, overflow: 'auto' } },
								}}
								placeholder="Please select"
								allowClear
								treeDefaultExpandAll
								treeData={apiTreeData}
								fieldNames={{ value: 'id', label: 'name' }}
							/>
						</Form.Item>
					</Col>
					<Col md={12}>
						<Form.Item
							label="Auth required"
							name={['auth_required']}
							rules={[
								{
									required: true,
								},
							]}
						>
							<Radio.Group
								options={yesNoOptions}
								block
								className="[&_.ant-radio-wrapper]:!justify-start"
							/>
						</Form.Item>
					</Col>
					<Col md={12}>
						<Form.Item
							label="Permission required"
							name={['permission_required']}
							rules={[
								{
									required: true,
								},
							]}
						>
							<Radio.Group
								options={yesNoOptions}
								block
								className="[&_.ant-radio-wrapper]:!justify-start"
							/>
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item
							label="Order"
							name={['sort_order']}
							rules={[
								{
									required: true,
								},
							]}
						>
							<InputNumber placeholder="Enter order" style={{ width: '100%' }} />
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item label="Description" name={['description']}>
							<TextArea
								placeholder="Enter description"
								autoSize={{ minRows: 4, maxRows: 4 }}
							/>
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</NewModal>
	);
};
export default ApiMainModal;
