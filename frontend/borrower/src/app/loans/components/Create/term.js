'use client';
import { Col, Form, Input, InputNumber, Row, Select } from 'antd';
import styles from './index.module.scss';
import { InputNumAmount, CommonFormDate, InputNameCombination } from '@repo/shared-methods';
import { useGetOptions } from '@/hooks/useGetOptions';
import { loan_active, loan_pending } from '@/enum/common';
import { capitalizeFirstLetter } from '@repo/shared-methods/utils';
const layoutCol = {
	xs: 24,
	sm: 24,
};
const CreateTerm = ({ stage, stepBtn, gutter, form, ...colProps }) => {
	const total_amount = Form.useWatch('total_amount', form);
	const { lienPositionOptions, prepaymentPenaltyOptions } = useGetOptions([
		{
			key: 'lien_position',
			valueName: 'lienPositionOptions',
		},
		{
			key: 'loan_prepayment_penalty',
			valueName: 'prepaymentPenaltyOptions',
		},
	]);

	const rateDom = (
		<Col {...layoutCol} {...(colProps || {})}>
			<Form.Item
				label="Rate"
				name="rate"
				rules={[{ required: true, message: 'Please input!' }]}
			>
				<InputNumber
					style={{ width: '100%' }}
					min={1}
					max={100}
					controls={false}
					placeholder="Enter Rate"
					suffix="%"
				/>
			</Form.Item>
		</Col>
	);

	const lenderDom = (
		<Col {...layoutCol} {...(colProps || {})}>
			<Form.Item
				label={stage === 'Pending' ? 'Target lender' : 'Lender'}
				name="lender_company"
				rules={[{ required: true, message: 'Please input!' }]}
				getValueFromEvent={(v) => capitalizeFirstLetter(v.target.value)}
			>
				<Input placeholder={`Enter ${stage === 'Pending' ? 'target ' : ''}lender`} />
			</Form.Item>
		</Col>
	);
	return (
		<Row className={styles.create} gutter={gutter}>
			<Form.List name="applicants">
				{(fields) => {
					return (
						<>
							{fields?.map((field, index) => (
								<InputNameCombination
									key={field.key}
									layoutCol={{ ...layoutCol, md: 8 }}
									firstName={[field.name, 'first_name']}
									middleName={[field.name, 'mid_name']}
									lastName={[field.name, 'last_name']}
									title={`Borrower ${index + 1}`}
									config={{
										firstNameRequired: !index,
										lastNameRequired: !index,
									}}
									listName="applicants"
									form={form}
									className="!ml-0 !mr-0"
								/>
							))}
						</>
					);
				}}
			</Form.List>
			{stage === loan_pending && (
				<>
					{lenderDom}
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Target closing date"
							name="closing_date"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<CommonFormDate placeholder="Select target closing date" />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Holdback"
							name="holdback_amount"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumAmount placeholder="Enter holdback amount" />
						</Form.Item>
					</Col>
					{rateDom}
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Estimated closing fund"
							name="closing_fund"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumAmount
								placeholder="Enter Estimated closing fund"
								max={total_amount}
							/>
						</Form.Item>
					</Col>
				</>
			)}
			{stage !== loan_pending && (
				<>
					{lenderDom}
					{rateDom}
				</>
			)}
			{stage === loan_active && (
				<>
					{/* <Col {...layoutCol} {...(colProps||{})}>
						<Form.Item
							label="Point"
							name="point"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={1}
								max={10}
								controls={false}
								placeholder="Enter point"
							/>
						</Form.Item>
					</Col> */}
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Lien position"
							name="lien_position"
							rules={[{ required: true, message: 'Please select Lien position!' }]}
						>
							<Select
								options={lienPositionOptions}
								placeholder="Select Lien position"
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Maturity date"
							name="maturity_date"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<CommonFormDate placeholder="Select Maturity date" />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Prepayment penalty"
							name="prepayment_penalty"
							rules={[{ required: true, message: 'Please input!' }]}
							getValueFromEvent={(e) => Number(e)}
						>
							<Select
								options={prepaymentPenaltyOptions}
								placeholder="Enter Prepayment penalty"
								suffixIcon={<div className="text-sm">Months</div>}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item
							label="Servicing company"
							name="service_company"
							rules={[{ required: true, message: 'Please input!' }]}
						>
							<Input placeholder="Enter Servicing company" />
						</Form.Item>
					</Col>
				</>
			)}
			{stage === loan_pending && (
				<>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Original contract date" name="original_contract_date">
							<CommonFormDate />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Extension date" name="extension_date">
							<CommonFormDate />
						</Form.Item>
					</Col>
					<Col {...layoutCol} {...(colProps || {})}>
						<Form.Item label="Loan processor" name="loan_processor">
							<Input placeholder="Enter Loan processor" />
						</Form.Item>
					</Col>
				</>
			)}
			{stepBtn}
		</Row>
	);
};
export default CreateTerm;
