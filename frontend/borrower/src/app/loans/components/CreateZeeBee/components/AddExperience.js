'use client';
import React, {useState} from 'react';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import {Col, Form, Input, Radio, Row} from 'antd';
import {useGetOptions} from '@repo/shared-methods/hooks';
import {CommonFormDate, InputNumAmount, NewModal} from '@repo/shared-methods';
import {yesNoOptions} from '@/enum/common';

export const AddExperience = ({ title, rowData = {}, type, onOk, onCancel }) => {
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const layoutCol = {
		md: 12,
	};
	const { getSelectFormDom } = useGetOptions([
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);
	return (
		<NewModal
			width={572}
			title={title}
			onCancel={onCancel}
			okText={title}
			confirmLoading={confirmLoading}
			onOk={async () => {
				try {
					setConfirmLoading(true);
					const values = await form.validateFields();
					if (rowData?.id) values.id = rowData?.id;
					await onOk(values, rowData?.id ? 'Edit' : 'Add');
					setConfirmLoading(false);
					onCancel();
				} catch {
					setConfirmLoading(false);
				}
			}}
		>
			<Form
				form={form}
				name="addExperienceForm"
				layout="vertical"
				autoComplete="off"
				className="w-full"
				size="large"
				disabled={confirmLoading}
				initialValues={{
					...rowData,
				}}
			>
				<Row gutter={24}>
					<AddressForm
						addressName={['address', 'address1']}
						cityName={['address', 'city']}
						stateName={['address', 'state']}
						addressLabel="Address"
						isShowAddress2={false}
						isShowZip={false}
						form={form}
						stateCol={layoutCol}
					/>
					<Col {...layoutCol}>
						<Form.Item
							label="Ownership on title" // 所有权人
							name="ownership_on_title"
							rules={[
								{
									required: true,
									message: 'Please enter ownership on title',
								},
							]}
						>
							<Input placeholder="Enter owner name" />
						</Form.Item>
					</Col>
					{getSelectFormDom(
						{
							label: 'Property type', // 房产类型
							name: 'property_type',
							placeholder: 'Select property type',
							valueName: 'propertyTypeOptions',
							requiredMessage: 'Please select property type!',
						},
						layoutCol,
					)}
					<Col {...layoutCol}>
						<Form.Item
							label="Purchase price" // 购买价格
							name="purchase_amount"
							rules={[
								{
									required: true,
									message: 'Please enter purchase price!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter purchase price" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Purchase date" // 购买时间
							name="purchase_date"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<CommonFormDate />
						</Form.Item>
					</Col>
					{type === 'transactionHistory' && (
						<Col {...layoutCol}>
							<Form.Item
								label="Rehab budget /Construction cost" // 重建预算
								name="budget_or_cost"
								rules={[
									{
										required: true,
										message: 'Please enter rehab budget/Construction cost!',
									},
								]}
							>
								<InputNumAmount placeholder="Enter rehab budget/Construction cost" />
							</Form.Item>
						</Col>
					)}
					{type === 'realEstate' && (
						<Col {...layoutCol}>
							<Form.Item
								label="Tenant Occupied" // 租户已入驻
								name="is_tenant_occupied"
								rules={[
									{
										required: true,
										message: 'Please select tenant Occupied!',
									},
								]}
							>
								<Radio.Group options={yesNoOptions} />
							</Form.Item>
						</Col>
					)}
					<Col {...layoutCol}>
						<Form.Item
							label="Sale price" // 售价
							name="sold_amount"
							rules={[
								{
									required: true,
									message: 'Please enter sale price!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter sale price amount" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Date Sold" // 售出时间
							name="sold_date"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<CommonFormDate />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Profit/ Loss" // 利润
							name="profit"
							rules={[
								{
									required: true,
									message: 'Please enter profit/loss!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter profit/loss amount" />
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</NewModal>
	);
};
