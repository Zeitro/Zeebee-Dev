'use client';
import { UploadOutlined } from '@ant-design/icons';
import {
	AddressForm,
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
import { Button, Col, Form, Input, InputNumber, Row, Image, message } from 'antd';
import { useState, forwardRef, useImperativeHandle, useEffect, useRef } from 'react';
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
		children,
		passwordName = 'password',
		confirmPasswordName = 'confirm_password',
		firstName = 'first_name',
		lastName = 'last_name',
		orgName = 'name',
		emailName = 'email',
		contactName = 'contact_name',
		contactEmailName = 'contact_email',
		contactPhoneName = 'contact_phone',
		taglineName = 'tagline',
		websiteName = 'website',
		loanAmountMinName = ['metadata', 'loan_amount_range', 'min'],
		loanAmountMaxName = ['metadata', 'loan_amount_range', 'max'],
		loanRateMinName = ['metadata', 'loan_rate_range', 'min'],
		loanRateMaxName = ['metadata', 'loan_rate_range', 'max'],
		loanPointMinName = ['metadata', 'loan_point_range', 'min'],
		loanPointMaxName = ['metadata', 'loan_point_range', 'max'],
		loanTermMinName = ['metadata', 'loan_term_range', 'min'],
		loanTermMaxName = ['metadata', 'loan_term_range', 'max'],
		logoName = 'logo',
		addressName = ['address', 'address1'],
		cityName = ['address', 'city'],
		stateName = ['address', 'state'],
		zipName = ['address', 'zip_code'],
	},
	ref,
) => {
	const isInit = useRef(true);
	const uploadRef = useRef(null);
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [fileList, setFileList] = useState([]);
	const to = <div className="flex items-center justify-center px-4 flex-shrink-0">to</div>;

	const handleSaveOrgProfile = async () => {
		try {
			if (!api) return;
			const values = await form.validateFields();
			let urls = [];
			const oldLogo = _.get(rowData, logoName);
			if (rowData?.id && oldLogo) {
				urls = [oldLogo];
			}
			if (fileList?.length) {
				urls = await uploadRef?.current?.triggerUpload();
			}
			if (!urls?.[0]) {
				message.open({
					content: 'Please select a logo',
					type: 'error',
				});
				return Promise.reject('logo is required');
			}
			_.set(values, logoName, urls?.[0]);
			const website = _.get(values, websiteName);
			if (website) {
				_.set(values, websiteName, `https://${website}`);
			}
			setConfirmLoading(true);
			await request(api, {
				method,
				body: values,
			});
			setConfirmLoading(false);
			successCallback?.();
		} catch (error) {
			setConfirmLoading(false);
			throw error?.message;
		}
	};

	useEffect(() => {
		if (rowData?.id && isInit?.current) {
			const info = _.cloneDeep(rowData);
			const website = _.get(rowData, websiteName)?.replace(/^https?:\/\//, '');
			_.set(info, websiteName, website || '');
			form.setFieldsValue(info);
			isInit.current = false;
		}
	}, [form, rowData, websiteName]);

	useImperativeHandle(ref, () => ({
		handleSaveOrgProfile,
		getInfo: () => {
			return form.getFieldsValue();
		},
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
			className={className}
		>
			<Row gutter={24}>
				<Col {...leftCol} className="!flex flex-col items-center overflow-hidden">
					<div className="w-[104px] h-[104px] border border-base-border rounded-xs px-1 py-1 mb-4">
						<Image
							src={fileList?.[0]?.thumbUrl || _.get(rowData, logoName) || null}
							alt=""
							fallback="/images/noImg.png"
							width={'100%'}
							height={'100%'}
							preview={!!fileList?.[0]?.thumbUrl || !!_.get(rowData, logoName)}
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
							<div className="text-2xl font-medium text-black-dark pt-6 break-words w-full overflow-hidden text-center">
								{_.get(rowData, orgName)}
							</div>
							<div className="text-xl font-medium text-[rgba(0,0,0,0.65)] pt-3 break-words w-full overflow-hidden text-center">
								{_.get(rowData, emailName)}
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
									firstName={firstName}
									lastName={lastName}
									hideMiddleName={true}
								/>
							</Row>
							<Form.Item
								label="Organization name"
								name={orgName}
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
								name={emailName}
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
							<ConfirmPassWord
								isShowRule={false}
								passwordName={passwordName}
								confirmPasswordName={confirmPasswordName}
							/>
						</>
					)}

					<Form.Item
						label="Contact name"
						name={contactName}
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
						name={contactEmailName}
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
						name={contactPhoneName}
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
					<AddressForm
						addressName={addressName}
						cityName={cityName}
						stateName={stateName}
						zipName={zipName}
						stateCol={{ xs: 24, sm: 24, md: 8 }}
						cityCol={{ xs: 24, sm: 24, md: 8 }}
						zipCol={{ xs: 24, sm: 24, md: 8 }}
						isHaveCol={false}
						isShowAddress2={false}
						isShowRow={true}
					/>
					<Form.Item label="Organization tagline" name={taglineName}>
						<TextArea
							placeholder="Organization tagline"
							autoSize={{ minRows: 2, maxRows: 2 }}
						/>
					</Form.Item>
					<Form.Item
						label="Website link"
						name={websiteName}
						getValueFromEvent={(v) => v.target.value?.replace(/^https?:\/\//, '')}
					>
						<Input addonBefore="https://" placeholder="Website link" allowClear />
					</Form.Item>
					<div className="!flex">
						<Form.Item
							label="Loan amount (US $)"
							name={loanAmountMinName}
							className="flex-1"
						>
							<InputNumAmount placeholder="Min.loan amount" className="!w-full" />
						</Form.Item>
						{to}
						<Form.Item label=" " name={loanAmountMaxName} className="flex-1">
							<InputNumAmount placeholder="Max.loan amount" className="!w-full" />
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item
							label="Loan term (Months)"
							name={loanTermMinName}
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
						<Form.Item label=" " name={loanTermMaxName} className="flex-1">
							<InputNumber
								placeholder="Max.loan term"
								className="!w-full"
								controls={false}
								// suffix="Months"
							/>
						</Form.Item>
					</div>
					<div className="!flex">
						<Form.Item label="Loan rate" name={loanRateMinName} className="flex-1">
							<InputNumber
								placeholder="Min.loan rate"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={loanRateMaxName} className="flex-1">
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
						<Form.Item label="Loan points" name={loanPointMinName} className="flex-1">
							<InputNumber
								placeholder="Min.loan points"
								className="!w-full"
								controls={false}
								suffix="%"
								max={100}
							/>
						</Form.Item>
						{to}
						<Form.Item label=" " name={loanPointMaxName} className="flex-1">
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
