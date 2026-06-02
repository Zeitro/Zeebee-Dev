import React, { useState } from 'react';
import { Row, Col, Form, Select, InputNumber } from 'antd';
// import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import { useGetOptions } from '@repo/shared-methods/hooks';
import { layoutCol } from '@repo/shared-methods/enum';
import {
	livingAreaProps,
	changeUnitByPropertyType,
	createRequest,
} from '@repo/shared-methods/utils';
import { NewModal, AddressForm } from '@repo/shared-methods';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const AddModal = ({ setOpen, rowData, handleSearch }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const [confirmLoading, setConfirmLoading] = useState(false);
	// const [loading, setLoading] = useState(false);
	const [unitInfo, setUnitInfo] = useState({});
	const { propertyTypeOptions } = useGetOptions([
		{
			key: 'property_type',
			valueName: 'propertyTypeOptions',
		},
	]);
	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			const value = await form.validateFields();
			let url = '/property';
			let method = 'POST';
			if (rowData?.id) {
				url = `/property/${rowData?.id}`;
				method = 'PUT';
			}
			const res = await request(url, {
				method,
				body: { ...value },
			});
			console.log(res);

			handleSearch?.(res?.id);
			setConfirmLoading(false);
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	useDeferredSyncEffect({
		predicate: () => rowData?.id,
		deferred: () => {
			const info = changeUnitByPropertyType(rowData?.property_type, rowData?.unit);
			setUnitInfo({ ...info });
			form.setFieldsValue({
				...rowData,
			});
		},
		deps: [form, rowData],
		delay: 0,
	});

	return (
		<NewModal
			width={700}
			title={`${rowData?.id ? 'Edit' : 'Add'} a property`}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				className="h-[60vh] overflow-hidden overflow-y-auto"
				size="large"
				disabled={confirmLoading}
				scrollToFirstError
			>
				<Row gutter={16}>
					<AddressForm
						addressName={['address', 'address1']}
						addressName2={['address', 'address2']}
						cityName={['address', 'city']}
						stateName={['address', 'state']}
						zipName={['address', 'zip_code']}
						// addressLabel="Subject property address 1"
						// addressLabel2="Subject property address 2"
						form={form}
						stateCol={{ ...layoutCol, md: 8 }}
						cityCol={{ ...layoutCol, md: 8 }}
						zipCol={{ ...layoutCol, md: 8 }}
						disabled={rowData?.id}
					/>
					<Col {...layoutCol}>
						<Form.Item
							label="Property type"
							name={['property_type']}
							rules={[{ required: true, message: 'Please select property type!' }]}
						>
							<Select
								options={propertyTypeOptions}
								placeholder="Select property type"
								onChange={(e) => {
									const info = changeUnitByPropertyType(e);
									setUnitInfo({ ...info });
									form.setFieldValue(['unit'], info.unit);
								}}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Living area sqft"
							name={['living_area_sqft']}
							rules={[
								{
									required: true,
									message: 'Please enter Living area sqft!',
								},
							]}
						>
							<InputNumber
								style={{ width: '100%' }}
								{...livingAreaProps}
								placeholder="Enter Living area sqft"
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Number of units"
							name={['unit']}
							rules={[
								{
									required: true,
									message: 'Please enter Number of units!',
								},
							]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={unitInfo?.minUnit}
								max={unitInfo?.maxUnit}
								controls={false}
								placeholder="Enter number of units"
								disabled={unitInfo?.disabled}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Lot area sqft"
							name={['lot_area_sqft']}
							rules={[
								{
									required: true,
									message: 'Please enter Lot area sqft!',
								},
							]}
						>
							<InputNumber
								style={{ width: '100%' }}
								{...livingAreaProps}
								placeholder="Enter Living area sqft"
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Number of bedrooms"
							name={['bedrooms']}
							rules={[
								{
									required: true,
									message: 'Please enter Number of bedrooms!',
								},
							]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={0}
								max={100}
								controls={false}
								placeholder="Enter Number of bedrooms"
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Number of bathrooms"
							name={['bathrooms']}
							rules={[
								{
									required: true,
									message: 'Please enter Number of bathrooms!',
								},
							]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={0}
								max={100}
								controls={false}
								placeholder="Enter Number of bathrooms"
							/>
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</NewModal>
	);
};
export default AddModal;
