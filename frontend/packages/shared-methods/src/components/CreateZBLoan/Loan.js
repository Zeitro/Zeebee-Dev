'use client';
import React, { useEffect } from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Radio, Input, InputNumber } from 'antd';
import { CommonFormDate } from '../CommonFormDate';
import { InputNumAmount } from '../InputNumAmount';
import { useGetOptions } from '@repo/shared-methods/hooks';
import { layoutCol, titleCol, yesNoOptions } from '@repo/shared-methods/enum';
import { maxTerm, maxAmount } from '../../lib/format';

export const LoanDetail = ({ gutter = 24, form }) => {
	const loanPurpose = Form.useWatch('loan_purpose', form);
	const loanType = Form.useWatch('loan_type', form) || 'Rental';
	const total_amount = Form.useWatch('total_amount', form) || maxAmount;

	const est_yearly_rental_income_amount = Form.useWatch(
		['loan_property', 'est_yearly_rental_income_amount'],
		form,
	);
	const yearly_liability_expense_amount = Form.useWatch(
		['loan_property', 'yearly_liability_expense_amount'],
		form,
	);
	const has_assignment_fee = Form.useWatch(['loan_property', 'has_assignment_fee'], form);
	const has_subordinate_debt = Form.useWatch(['loan_property', 'has_subordinate_debt'], form);
	const isPurchase = loanPurpose === 'Purchase';

	const { getSelectFormDom } = useGetOptions([
		{
			key: 'business_strategy',
			valueName: 'businessStrategyOptions',
		},
		{
			key: 'loan_prepayment_penalty',
			valueName: 'prepaymentPenaltyOptions',
		},
	]);

	useEffect(() => {
		if (est_yearly_rental_income_amount && yearly_liability_expense_amount) {
			const num = est_yearly_rental_income_amount / yearly_liability_expense_amount;
			form.setFieldsValue({ dscr: Math.round(num * 100) / 100 });
		}
	}, [est_yearly_rental_income_amount, form, yearly_liability_expense_amount]);
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...titleCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Loan detail
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Purchase price"
					name={['loan_property', 'purchase_amount']}
					rules={[{ required: true, message: 'Please select purchase price!' }]}
				>
					<InputNumAmount placeholder="Enter purchase price" />
				</Form.Item>
			</Col>
			{isPurchase && (
				<>
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Assignment fee"
							name={['loan_property', 'has_assignment_fee']}
							rules={[
								{
									required: true,
									message: 'Please enter assignment fee!',
								},
							]}
						>
							<Radio.Group
								options={yesNoOptions}
								style={{ width: '100%' }}
								className={styles.yesNo}
							/>
						</Form.Item>
					</Col>
					{has_assignment_fee && (
						<Col {...layoutCol} md={12}>
							<Form.Item
								label="Assignment fee"
								name={['loan_property', 'assignment_fee_amount']}
								rules={[
									{
										required: true,
										message: 'Please enter assignment fee amount!',
									},
								]}
							>
								<InputNumAmount placeholder="Enter assignment fee amount" />
							</Form.Item>
						</Col>
					)}
				</>
			)}

			<Col {...layoutCol}>
				{isPurchase ? (
					<Form.Item
						label="Scheduled close of escrow date" // 预定房产交易日
						name={['loan_property', 'scheduled_close_escrow_date']}
						rules={[
							{
								required: true,
								message: 'Please select scheduled close of escrow date!',
							},
						]}
					>
						<CommonFormDate placeholder="Enter scheduled close of escrow date" />
					</Form.Item>
				) : (
					<Form.Item
						label="Date of purchase"
						name={['loan_property', 'purchase_date']}
						rules={[
							{
								required: true,
								message: 'Please select date of purchase!',
							},
						]}
					>
						<CommonFormDate placeholder="Enter date of purchase" />
					</Form.Item>
				)}
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Initial loan" // 初始贷款金额（即第一笔发放的资金）
					name="initial_amount"
					rules={[
						{
							required: true,
							message: 'Please enter initial loan amount!',
						},
					]}
				>
					<InputNumAmount placeholder="Enter initial loan amount" max={total_amount} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Holdback" // 贷款中暂时保留未发放的金额
					name="holdback_amount"
					rules={[
						{
							required: true,
							message: 'Please enter holdback amount!',
						},
					]}
				>
					<InputNumAmount placeholder="Enter holdback amount" max={total_amount} />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Total loan amount" // 总贷款金额
					name="total_amount"
					rules={[
						{
							required: true,
							message: 'Please enter total loan amount!',
						},
					]}
				>
					<InputNumAmount placeholder="Enter total loan amount" />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Loan term requested" // 借款期限
					name="loan_term"
					rules={[{ required: true, message: 'Please enter loan term requested!' }]}
				>
					<InputNumber
						min={1}
						max={maxTerm}
						controls={false}
						suffix="Months"
						style={{ width: '100%' }}
						placeholder="Enter loan term requested"
					/>
				</Form.Item>
			</Col>
			{!isPurchase && (
				<Col {...layoutCol}>
					<Form.Item
						label="Estimated Current Value (LTV)" // 预计当前估值
						name={['loan_property', 'est_property_value_amount']}
						rules={[
							{
								required: true,
								message: 'Please enter estimated LTV!',
							},
						]}
					>
						<InputNumAmount placeholder="Enter estimated LTV" />
					</Form.Item>
				</Col>
			)}
			<Col {...layoutCol}>
				<Form.Item
					label="Estimated after repair value (ARV)" // 修复后预计估值
					name={['loan_property', 'arv']}
					rules={[
						{
							required: true,
							message: 'Please enter estimated ARV!',
						},
					]}
				>
					<InputNumAmount placeholder="Enter estimated ARV" />
				</Form.Item>
			</Col>
			{!isPurchase && (
				<>
					<Col {...layoutCol}>
						<Form.Item
							label="Outstanding loan amount" // 当前尚未偿还的贷款本金
							name="principal_amount"
							rules={[
								{
									required: true,
									message: 'Please enter outstanding Loan Amount!',
								},
							]}
						>
							<InputNumAmount
								placeholder="Enter outstanding Loan Amount"
								max={total_amount}
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Existing Lender" // 现有贷款方，出借人（房子抵押贷款前购买时的出借人）
							name={['loan_property', 'existing_lender_company']}
							rules={[{ required: true, message: 'Please enter existing Lender!' }]}
						>
							<Input style={{ width: '100%' }} placeholder="Enter existing Lender" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Existing loan rate" // 房子抵押贷款前购买时的贷款利率
							name={['existing_rate', 'existing_rate']}
							rules={[
								{ required: true, message: 'Please enter existing loan rate!' },
							]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={1}
								max={100}
								controls={false}
								placeholder="Enter existing loan rate"
								suffix="%"
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Existing loan Maturity date" // 现有贷款到期日
							name={['loan_property', 'existing_maturity_date']}
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<CommonFormDate placeholder="Enter existing loan maturity date" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Subordinate debt"
							name={['loan_property', 'has_subordinate_debt']}
							rules={[
								{
									required: true,
									message: 'Please select subordinate debt!',
								},
							]}
						>
							<Radio.Group options={yesNoOptions} className={styles.yesNo} />
						</Form.Item>
					</Col>
					{has_subordinate_debt && (
						<Col {...layoutCol}>
							<Form.Item
								label="Subordinate debt"
								name={['loan_property', 'subordinate_debt_amount']}
								rules={[
									{
										required: true,
										message: 'Please select subordinate debt amount!',
									},
								]}
							>
								<InputNumAmount placeholder="Enter subordinate debt amount" />
							</Form.Item>
						</Col>
					)}
				</>
			)}
			<Col {...layoutCol}>
				<Form.Item
					label="Liquidity" // 流动资金
					name={['loan_property', 'liquidity_amount']}
					rules={[
						{
							required: true,
							message: 'Please enter liquidity!',
						},
					]}
				>
					<InputNumAmount placeholder="Enter liquidity" />
				</Form.Item>
			</Col>
			{getSelectFormDom(
				{
					label: 'Business strategy', // 商业策略
					name: 'business_strategy',
					mode: 'multiple',
					placeholder: 'Select business strategy',
					valueName: 'businessStrategyOptions',
					requiredMessage: 'Please enter business strategy!',
				},
				layoutCol,
			)}
			{!isPurchase && (
				<Col {...layoutCol}>
					<Form.Item
						label="Mid construction" // 中建
						name={['loan_property', 'mid_construction']}
						rules={[
							{
								required: true,
								message: 'Please enter mid construction!',
							},
						]}
					>
						<Radio.Group options={yesNoOptions} className={styles.yesNo} />
					</Form.Item>
				</Col>
			)}
			{['Rehab', 'Ground up construction'].includes(loanType) && (
				<>
					<Col {...titleCol} className="!text-xl font-medium text-[#000000a6] pb-6">
						Budget detail
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Original budget" // 最初预算
							name={['loan_property', 'budget', 'original_budget']}
							rules={[
								{
									required: true,
									message: 'Please enter original budget!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter original budget" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Remaining budget" // 剩余预算
							name={['loan_property', 'budget', 'remaining_budget']}
							rules={[
								{
									required: true,
									message: 'Please enter remaining budget!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter remaining budget" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Cost spent to date" // 截至目前所用花费
							name={['loan_property', 'budget', 'cost_spent_to_date_amount']}
							rules={[
								{
									required: true,
									message: 'Please enter cost spent to date!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter cost spent to date" />
						</Form.Item>
					</Col>
				</>
			)}
			{loanType === 'Rental' && (
				<>
					<Col {...titleCol} className="!text-xl font-medium text-[#000000a6] pb-6">
						{/* Rental detail */}
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Estimated net yearly rental income" // 预计年租金净收入
							name={['loan_property', 'est_yearly_rental_income_amount']}
							rules={[
								{
									required: true,
									message: 'Please enter income amount!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter income amount" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Yearly liability expenses" // 年度负债支出
							name={['loan_property', 'yearly_liability_expense_amount']}
							rules={[
								{
									required: true,
									message: 'Please enter yearly liability expenses!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter yearly liability expenses" />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="DSCR" // 债务覆盖率，衡量企业偿债能力的关键指标
							name={['loan_property', 'dscr']}
							rules={[{ required: true, message: 'Please enter DSCR!' }]}
						>
							<InputNumber
								style={{ width: '100%' }}
								min={0}
								controls={false}
								placeholder="Enter DSCR"
								suffix="%"
								disabled
							/>
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Reserves" // 储备金
							name="reserves_amount"
							rules={[
								{
									required: true,
									message: 'Please enter reserves!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter reserves" />
						</Form.Item>
					</Col>
					{getSelectFormDom(
						{
							label: 'Prepayment penalty', // 提前还款罚金
							name: 'prepayment_penalty',
							valueName: 'prepaymentPenaltyOptions',
							requiredMessage: 'Please enter Prepayment penalty!',
							suffixIcon: <div className="text-sm">Months</div>,
							formProps: {
								getValueFromEvent: (e) => Number(e),
							},
						},
						layoutCol,
					)}
					<Col {...layoutCol}></Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Delayed financing" // 延期融资
							name={['loan_property', 'has_delayed_financing']}
							rules={[
								{
									required: true,
									message: 'Please select!',
								},
							]}
						>
							<Radio.Group options={yesNoOptions} className={styles.yesNo} />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Collateralized/Blanket" // 抵押/保函
							name={['loan_property', 'has_collateralized_blanket']}
							rules={[
								{
									required: true,
									message: 'Please select!',
								},
							]}
						>
							<Radio.Group options={yesNoOptions} className={styles.yesNo} />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Short term rental" // 短期租赁
							name={['loan_property', 'has_short_term_rental']}
							rules={[
								{
									required: true,
									message: 'Please select!',
								},
							]}
						>
							<Radio.Group options={yesNoOptions} className={styles.yesNo} />
						</Form.Item>
					</Col>
					<Col {...layoutCol}>
						<Form.Item
							label="Prosperities vacant" // 资产空置
							name="has_properties_vacant"
							rules={[
								{
									required: true,
									message: 'Please select!',
								},
							]}
						>
							<Radio.Group options={yesNoOptions} className={styles.yesNo} />
						</Form.Item>
					</Col>
				</>
			)}
		</Row>
	);
};
export default LoanDetail;
