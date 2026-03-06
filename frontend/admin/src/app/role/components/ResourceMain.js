import React, { useEffect, useState } from 'react';
import { Row, Col, Form, Select, Input, message, Radio, TreeSelect, InputNumber, Tag } from 'antd';
import { NewModal } from '@repo/shared-methods';
import { yesNoOptions, systemTypeOptions, titleCol } from '@repo/shared-methods/enum';
import { createRequest, eventBus } from '@repo/shared-methods/utils';
import useData from './useData';
const { TextArea } = Input;
const MAX_COUNT = 5;
const ResourceMain = ({ setOpen, rowData, eventBusPrefix, api }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const resource_type = Form.useWatch('resource_type', form);
	const api_ids = Form.useWatch('api_ids', form);
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const [record, setRecord] = useState({});
	const [treeData, setTreeData] = useState([]);
	const { getApiTree, apiTreeData } = useData();

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/resource';
			let method = 'POST';
			let msg = 'Api creation successful.';
			if (rowData?.id) {
				// value.id = rowData?.id;
				url = `/resource/${rowData?.id}`;
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
		} catch (err) {
			console.log(err);

			setConfirmLoading(false);
		}
	};

	const getDetails = async (type) => {
		try {
			setLoading(true);
			const payload = new URLSearchParams({ type });
			const url = `/resource/${rowData?.id}?${payload}`;
			const res = await request(url);
			res.api_ids = res?.apis?.map((item) => item.id);
			setRecord(res);
			form.setFieldsValue(res);
			setLoading(false);
		} catch (error) {
			setLoading(false);
			throw error;
		}
	};

	const getList = async (v) => {
		try {
			setLoading(true);
			const payload = new URLSearchParams({ type: v });
			const url = `${api}?${payload}`;
			const res = await request(url);
			setTreeData(res?.list);
			setLoading(false);
		} catch (error) {
			setLoading(false);
			throw error;
		}
	};

	useEffect(() => {
		if (rowData?.type && api) {
			const getFn = async () => {
				try {
					setLoading(true);
					await Promise.all([getList(rowData?.type), getApiTree(rowData?.type)]);
					if (rowData?.id) {
						await getDetails(rowData?.type);
					}
					setLoading(false);
				} catch (error) {
					setLoading(false);
					throw error;
				}
			};
			getFn();
		}
	}, [rowData?.type, api]);

	return (
		<NewModal
			width={700}
			title={`${rowData?.id ? 'Edit' : 'Add'} resource`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				form={form}
				name="resource-form"
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
					enable_data_scope: false,
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
						<Form.Item label="Resource Type" name={['resource_type']}>
							<Select
								options={[
									{ value: 'Menu', label: 'Menu' },
									{ value: 'Button', label: 'Button' },
								]}
								placeholder="Select"
							/>
						</Form.Item>
					</Col>
					<Col md={24}>
						<Form.Item label="Path" name={['path']}>
							<Input placeholder="Enter path" />
						</Form.Item>
					</Col>
					<Col md={12}>
						<Form.Item
							label="Tag"
							name={['tag']}
							rules={[
								{
									required: true,
								},
							]}
						>
							<Input placeholder="Enter tag" />
						</Form.Item>
					</Col>

					<Col md={12}>
						<Form.Item label="Icon" name={['icon']}>
							<Input placeholder="Enter icon" />
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
									getList(v);
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
						<Form.Item label="Parent resource" name={['parent_id']}>
							<TreeSelect
								showSearch
								style={{ width: '100%' }}
								styles={{
									popup: { root: { maxHeight: 400, overflow: 'auto' } },
								}}
								placeholder="Please select"
								allowClear
								treeDefaultExpandAll
								treeData={treeData}
								fieldNames={{ value: 'id', label: 'name' }}
							/>
						</Form.Item>
					</Col>
					{resource_type === 'Menu' && (
						<Col md={12}>
							<Form.Item
								label="Show menu"
								name={['is_show']}
								rules={[
									{
										required: true,
									},
								]}
								extra="Yes: Show this menu in the sidebar / No: Hide."
							>
								<Radio.Group
									options={yesNoOptions}
									block
									className="[&_.ant-radio-wrapper]:!justify-start"
								/>
							</Form.Item>
						</Col>
					)}
					<Col md={12}>
						<Form.Item
							label="Configurable data permissions"
							name={['enable_data_scope']}
							rules={[
								{
									required: true,
								},
							]}
							extra="Yes: Show,The system can be configured by itself.  No: Hidden configuration"
						>
							<Radio.Group
								options={yesNoOptions}
								block
								className="[&_.ant-radio-wrapper]:!justify-start"
							/>
						</Form.Item>
					</Col>
					<Col md={12}>
						<Form.Item label="Order" name={['sort_order']} rules={[{ required: true }]}>
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
					<Col {...titleCol}>
						<Form.Item
							label="Related APIs"
							name={['api_ids']}
							rules={[{ required: true }]}
						>
							<TreeSelect
								showSearch
								multiple
								treeCheckable
								style={{ width: '100%' }}
								placeholder="Please select"
								allowClear
								treeDefaultExpandAll
								treeData={apiTreeData}
								showCheckedStrategy={TreeSelect.SHOW_CHILD}
								maxTagCount={5}
								fieldNames={{ value: 'id', label: 'name' }}
								treeTitleRender={(item) => {
									return item?.path ? (
										<>
											{item.name} (
											<span className="text-base-label"> {item.path}</span>)
										</>
									) : (
										item?.name
									);
								}}
								tagRender={(item) => {
									return (
										<Tag
											closable={typeof item?.label !== 'string'}
											onClose={(removedTag) => {
												const newTags = api_ids.filter(
													(tag) => tag !== removedTag,
												);
												form.setFieldValue(newTags);
											}}
										>
											{item?.label?.props?.children?.[0] || item?.label}
										</Tag>
									);
								}}
							/>
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</NewModal>
	);
};
export default ResourceMain;
