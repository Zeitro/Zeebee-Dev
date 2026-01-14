'use client';
import { UploadOutlined } from '@ant-design/icons';
import {
	ConfirmPassWord,
	InputNameCombination,
	InputNumAmount,
	InputPhone,
	UploadFileButton,
} from '@repo/shared-methods';
import {
	capitalizeFirstLetter,
	emailRegex,
	phoneRegex,
	validateRule,
	createRequest,
} from '@repo/shared-methods/utils';
import { Button, Col, Form, Input, InputNumber, Row, Image } from 'antd';
import { useState, forwardRef, useImperativeHandle, useEffect, useCallback, useRef } from 'react';
import _ from 'lodash';
const { TextArea } = Input;
// 组织信息填写
const OrganizationProfile = (
	{
		api,
		method = 'PUT',
		rowData,
		successCallback,
		className,
		leftCol = { lg: 8, sm: 24, xl: 6 },
		rightCol = { lg: 16, sm: 24, xl: 12 },
		isShowAccountInput,
		metadataName = [],
		initMetadataName = [],
		children,
	},
	ref,
) => {
	const uploadRef = useRef(null);
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [fileList, setFileList] = useState([]);
	const isInit = useRef(true);
	const to = <div className="flex items-center justify-center px-4 flex-shrink-0">to</div>;

	const handleSaveOrgProfile = async () => {
		try {
			if (!api) return;
			const values = await form.validateFields();
			setConfirmLoading(true);
			let urls = [];
			if (fileList?.length) {
				urls = await uploadRef?.current?.triggerUpload();
			}
			await request(api, {
				method,
				body: {
					...values,
					logo: urls?.[0],
				},
			});
			setConfirmLoading(false);
			successCallback?.();
		} catch (error) {
			setConfirmLoading(false);
			throw error?.message;
		}
	};

	// 初始化赋值
	const setFormValue = useCallback(() => {
		if (rowData?.id) {
			const info = _.cloneDeep(rowData);
			if (initMetadataName?.length && isInit.current) {
				const v = _.get(info, initMetadataName);
				_.set(info, metadataName, v);
				form.setFieldsValue(info);
				isInit.current = false;
			}
		}
	}, [rowData, initMetadataName, metadataName, form]);

	useEffect(() => {
		setFormValue();
	}, [setFormValue]);

	useImperativeHandle(ref, () => ({
		handleSaveOrgProfile,
	}));
	return (
		<Form
			form={form}
			name="org-profile-form"
			layout="vertical"
			autoComplete="off"
			size="large"
			disabled={confirmLoading}
			scrollToFirstError
			initialValues={{ ...rowData }}
			className={className}
		>
			<Row gutter={24}>
				<Col {...leftCol} className="!flex flex-col items-center">
					<div className="w-[104px] h-[104px] border border-base-border rounded-xs px-1 py-1 mb-4">
						<Image
							src={fileList?.[0]?.thumbUrl || rowData?.logo || null}
							alt=""
							fallback="/images/noImg.png"
							width={'100%'}
							height={'100%'}
							preview={!!fileList?.[0]?.thumbUrl || !!rowData?.logo}
							className="object-contain"
						/>
					</div>
					<UploadFileButton
						ref={uploadRef}
						accept=".png,.jpg,.jpeg"
						showUploadList={false}
						immediateUpload={false}
						onGetUploadList={(fileList) => {
							setFileList(fileList);
						}}
					>
						<Button icon={<UploadOutlined />} size="middle">
							Upload new logo
						</Button>
					</UploadFileButton>

					{!isShowAccountInput && (
						<>
							<div className="text-2xl font-medium text-black-dark pt-6">
								{rowData?.name}
							</div>
							<div className="text-xl font-medium text-[rgba(0,0,0,0.65)] pt-3">
								{rowData?.owner?.email}
							</div>
						</>
					)}
				</Col>
				<Col {...rightCol}>
					{isShowAccountInput && (
						<>
							<Row gutter={16}>
								<InputNameCombination
									layoutCol={{ span: 12 }}
									firstName={['first_name']}
									lastName={['last_name']}
									hideMiddleName={true}
								/>
							</Row>
							<Form.Item
								label="Organization name"
								name={['org_name']}
								rules={[
									{
										required: true,
									},
								]}
							>
								<Input allowClear placeholder="Organization name" />
							</Form.Item>
							<Form.Item
								label="Email account"
								name={['email']}
								rules={[
									{
										required: true,
									},
									{
										validator: (_, value) =>
											validateRule(value, emailRegex, 'Invalid email format'),
									},
								]}
								extra="Email account used to log in to the ZeeBee lender system."
							>
								<Input placeholder="Email account" allowClear />
							</Form.Item>
							<ConfirmPassWord isShowRule={false} />
						</>
					)}

					<Form.Item
						label="Contact name"
						name={['contact_name']}
						rules={[
							{
								required: true,
							},
						]}
						getValueFromEvent={(e) => {
							return capitalizeFirstLetter(e.target.value);
						}}
					>
						<Input placeholder="Contact name" allowClear />
					</Form.Item>
					<Form.Item
						label="Contact email"
						name={['contact_email']}
						rules={[
							{
								required: true,
							},
							{
								validator: (_, value) =>
									validateRule(value, emailRegex, 'Invalid email format'),
							},
						]}
					>
						<Input placeholder="Contact email" allowClear />
					</Form.Item>
					<Form.Item
						label="Contact phone"
						name={['contact_phone']}
						rules={[
							{
								required: true,
							},
							{
								validator: (_, value) =>
									validateRule(value, phoneRegex, 'Invalid phone number'),
							},
						]}
					>
						<InputPhone placeholder="Contact phone" />
					</Form.Item>
					<Form.Item label="Organization tagline" name={['tagline']}>
						<TextArea
							placeholder="Organization tagline"
							autoSize={{ minRows: 2, maxRows: 2 }}
						/>
					</Form.Item>
					<Form.Item label="Website link" name={['website']}>
						<Input addonBefore="https://" placeholder="Website link" allowClear />
					</Form.Item>
					<div className="!flex">
						<Form.Item
							label="Loan amount (US $)"
							name={[...metadataName, 'loan_amount_range', 'min']}
							className="flex-1"
						>
							<InputNumAmount placeholder="Min.loan amount" className="!w-full" />
						</Form.Item>
						{to}
						<Form.Item
							label=" "
							name={[...metadataName, 'loan_amount_range', 'max']}
							className="flex-1"
						>
							<InputNumAmount placeholder="Max.loan amount" className="!w-full" />
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item
							label="Loan term (Months)"
							name={[...metadataName, 'loan_term_range', 'min']}
							className="flex-1"
						>
							<InputNumber
								placeholder="Min.loan term"
								className="!w-full"
								controls={false}
								// suffix="Months"
							/>
						</Form.Item>
						{to}
						<Form.Item
							label=" "
							name={[...metadataName, 'loan_term_range', 'max']}
							className="flex-1"
						>
							<InputNumber
								placeholder="Max.loan term"
								className="!w-full"
								controls={false}
								// suffix="Months"
							/>
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item
							label="Loan rate"
							name={[...metadataName, 'loan_rate_range', 'min']}
							className="flex-1"
						>
							<InputNumber
								placeholder="Min.loan rate"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
						{to}
						<Form.Item
							label=" "
							name={[...metadataName, 'loan_rate_range', 'max']}
							className="flex-1"
						>
							<InputNumber
								placeholder="Max.loan rate"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item
							label="Loan points"
							name={[...metadataName, 'loan_point_range', 'min']}
							className="flex-1"
						>
							<InputNumber
								placeholder="Min.loan points"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
						{to}
						<Form.Item
							label=" "
							name={[...metadataName, 'loan_point_range', 'max']}
							className="flex-1"
						>
							<InputNumber
								placeholder="Max.loan points"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
					</div>
					{children}
				</Col>
			</Row>
		</Form>
	);
};
OrganizationProfile.displayName = 'OrgProfile';
export const OrgProfile = forwardRef(OrganizationProfile);
