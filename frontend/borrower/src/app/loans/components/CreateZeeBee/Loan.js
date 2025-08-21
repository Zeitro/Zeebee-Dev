import React, { useEffect } from 'react';
import styles from './index.module.scss';
import { Row, Col, Form, Radio, Select, Input, InputNumber, DatePicker } from 'antd';
import AddressForm from '@/components/GoogleAutoComplete/AddressForm';
import InputPhone from '@/components/InputPhone';
import SelectTINType from '@/components/InputTIN/SelectTINType';
import CommonFormDate from '@/components/CommonDate';
import {
	creditScoreRegex,
	emailRegex,
	maxCreditScore,
	minCreditScore,
	phoneRegex,
} from '@/lib/format';
import InputNumAmount from '@/components/InputNumAmount';
import { yesNoOptions } from '@/enum/common';
import { formatStr } from '@/lib/date';
import { useGetOptions } from '@/hooks/useGetOptions';
const MAX_COUNT = 5;
const LoanDetail = ({ gutter = 24, form, addressProps }) => {
	const loanPurpose = Form.useWatch('loan_purpose', form);
	const loanType = Form.useWatch('loan_type', form) || 'Rental';
	const yearly_rental_income = Form.useWatch('yearly_rental_income', form);
	const yearly_liability_expenses = Form.useWatch('yearly_liability_expenses', form);
	const isPurchase = loanPurpose === 'purchase';
	const layoutCol = {
		xs: 24,
		sm: 24,
	};
	const { businessStrategyOptions, getSelectFormDom } = useGetOptions([
		{
			key: 'business_strategy',
			valueName: 'businessStrategyOptions',
		},
	]);

	useEffect(() => {
		if (yearly_rental_income && yearly_liability_expenses) {
			form.setFieldsValue({ DSCR: yearly_rental_income / yearly_liability_expenses });
		}
	}, [yearly_rental_income, yearly_liability_expenses]);
	return (
		<Row className={styles.createZeeBee} gutter={gutter}>
			<Col {...layoutCol} className="!text-2xl/8 font-medium text-[#000000a6] pb-8">
				Loan detail
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Purchase price"
					name="purchase_price"
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
							name="assignment_fee"
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
					<Col {...layoutCol} md={12}>
						<Form.Item
							label=" "
							name="assignment_fee"
							rules={[
								{
									required: true,
									message: 'Please enter assignment fee!',
								},
							]}
							className={styles.assignmentFee}
						>
							<InputNumAmount placeholder="Enter assignment fee" />
						</Form.Item>
					</Col>
				</>
			)}

			<Col {...layoutCol}>
				{isPurchase ? (
					<Form.Item
						label="Scheduled close of escrow date" // 预定房产交易日
						name="date_of_birth"
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
						name="date_of_purchase"
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
					name="initial_loan"
					rules={[
						{
							required: true,
							message: 'Please enter initial loan amount!',
						},
					]}
				>
					<InputNumAmount placeholder="Enter initial loan amount" />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Holdback" // 贷款中暂时保留未发放的金额
					name="holdback"
					rules={[
						{
							required: true,
							message: 'Please enter holdback amount!',
						},
					]}
				>
					<InputNumAmount placeholder="Enter holdback amount" />
				</Form.Item>
			</Col>
			<Col {...layoutCol}>
				<Form.Item
					label="Total loan amount" // 总贷款金额
					name="total_loan_amount"
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
					name="term"
					rules={[{ required: true, message: 'Please select loan term requested!' }]}
				>
					<InputNumber
						min={1}
						max={10000}
						controls={false}
						suffix="Months"
						style={{ width: '100%' }}
					/>
				</Form.Item>
			</Col>
			{!isPurchase && (
				<Col {...layoutCol}>
					<Form.Item
						label="Estimated Current Value (LTV)" // 预计当前估值
						name="total_loan_amount"
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
					name="total_loan_amount"
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
							name="outstanding_loan_Amount"
							rules={[
								{
									required: true,
									message: 'Please enter outstanding Loan Amount!',
								},
							]}
						>
							<InputNumAmount placeholder="Enter outstanding Loan Amount" />
						</Form.Item>
					</Col>
					{getSelectFormDom({
						label: 'Existing Lender', // 现有贷款方，出借人
						name: 'Existing Lender',
						placeholder: 'Select existing Lender',
						valueName: 'lenderOptions',
						requiredMessage: 'Please enter existing Lender!',
					})}
					<Col {...layoutCol}>
						<Form.Item
							label="Existing loan rate" // 现有贷款利率
							name="rate"
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
							name="maturity_date"
							rules={[{ required: true, message: 'Please select date!' }]}
						>
							<CommonFormDate placeholder="Enter existing loan maturity date" />
						</Form.Item>
					</Col>
				</>
			)}
			<Col {...layoutCol}>
				<Form.Item
					label="Liquidity" // 流动资金
					name="liquidity"
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
			{getSelectFormDom({
				label: 'Buidiness strategy', // 商业策略
				name: 'buidiness_strategy',
				mode: 'multiple',
				placeholder: 'Select buidiness strategy',
				valueName: 'businessStrategyOptions',
				requiredMessage: 'Please enter buidiness strategy!',
			})}
			{!isPurchase && (
				<Col {...layoutCol} md={12}>
					<Form.Item
						label="Mid construction" // 中建
						name="midConstruction"
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
			<Col {...layoutCol} className="!text-xl font-medium text-[#000000a6] pb-6">
				Budget detail
			</Col>
			{['Rehab', 'Ground up construction'].includes(loanType) && (
				<>
					<Col {...layoutCol}>
						<Form.Item
							label="Original budget" // 最初预算
							name="original_budget"
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
							name="remaining_budget"
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
							name="cost_spent_to_date"
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
					<Col {...layoutCol}>
						<Form.Item
							label="Est. Net yearly rental income" // 预计年租金收入
							name="yearly_rental_income"
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
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Yearly liability expenses" // 年度负债支出
							name="yearly_liability_expenses"
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
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="DSCR"
							name="DSCR"
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
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Reserves" // 储备金
							name="reserves"
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
						},
						{ md: 12 },
					)}
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Delayed financing" // 延期融资
							name="delayed_financing"
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
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Collateralized/Blanket" // 抵押/保函
							name="collateralized_blanket"
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
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Short term rental" // 短期租赁
							name="Short term rental"
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
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Prosperities vacant" // 资产空置
							name="Prosperities vacant"
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
