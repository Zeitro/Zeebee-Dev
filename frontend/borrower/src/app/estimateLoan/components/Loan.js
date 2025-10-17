import React, { useEffect } from 'react';
// import styles from './index.module.scss';
import { Row, Col, Form, InputNumber, Radio } from 'antd';
import InputNumAmount from '@/components/InputNumAmount';
import { layoutCol, titleCol } from '@/enum/common';
import { maxTerm_zeeBee, minTerm_zeeBee, validateAssignmentFeeAmount } from '@/lib/format';
import { yesNoOptions } from '@/enum/common';
import CommonFormDate from '@/components/CommonDate';
const LoanInfo = ({ gutter = 24, loan_purpose, loan_type }) => {
	const form = Form.useFormInstance();
	const has_assignment_fee = Form.useWatch('has_assignment_fee', form);
	const purchase_amount = Form.useWatch('purchase_amount', form);
	const total_amount = Form.useWatch('total_amount', form);
	const isPurchase = loan_purpose === 'Purchase';
	return (
		<Row gutter={gutter}>
			<Col {...titleCol} className="!text-3xl/8 font-medium pb-8">
				Loan information
			</Col>

			<Col {...titleCol}>
				<Form.Item
					label="Loan terms" // 借款期限
					name="loan_term"
					rules={[{ required: true }]}
				>
					<InputNumber
						min={minTerm_zeeBee}
						max={maxTerm_zeeBee}
						controls={false}
						suffix="Months"
						style={{ width: '100%' }}
						placeholder="Enter loan terms"
					/>
				</Form.Item>
			</Col>
			{!isPurchase && (
				<>
					<Col {...titleCol}>
						<Form.Item
							label="Exiting loan amount" //对应Outstanding loan amount
							name="principal_amount"
							rules={[{ required: true }]}
						>
							<InputNumAmount placeholder="Exiting loan amount" />
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item
							label="Purchase date"
							name="purchase_date"
							rules={[{ required: true }]}
						>
							<CommonFormDate placeholder="Enter date of purchase" />
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item
							label="Initial loan" // 初始贷款金额（即第一笔发放的资金）
							name="initial_amount"
							rules={[{ required: true }]}
						>
							<InputNumAmount
								placeholder="Enter initial loan amount"
								max={total_amount}
							/>
						</Form.Item>
					</Col>
				</>
			)}
			{['Rehab', 'Ground up construction'].includes(loan_type) && (
				<>
					<Col {...titleCol}>
						<Form.Item
							label={
								loan_type === 'Rehab'
									? 'Rehab loan amount'
									: 'Construction loan amount'
							} // 贷款中暂时保留未发放的金额-Zeebee中的Holdback
							name="holdback_amount"
							rules={[{ required: true }]}
						>
							<InputNumAmount
								placeholder={
									loan_type === 'Rehab'
										? 'Rehab loan amount'
										: 'Construction loan amount'
								}
								max={total_amount}
							/>
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item
							label={loan_type === 'Rehab' ? 'Rehab budget' : 'Construction budget'} // 最初预算-Zeebee字段Original budget
							name={['budget', 'original_budget']}
							rules={[{ required: true }]}
						>
							<InputNumAmount
								placeholder={
									loan_type === 'Rehab' ? 'Rehab budget' : 'Construction budget'
								}
							/>
						</Form.Item>
					</Col>
					<Col {...titleCol}>
						<Form.Item
							label="Cost spent to date" // 截至目前所用花费
							name={['budget', 'cost_spent_to_date_amount']}
							rules={[{ required: true }]}
						>
							<InputNumAmount placeholder="Cost spent to date" />
						</Form.Item>
					</Col>
				</>
			)}
			{loan_type === 'Rental' && (
				<>
					<Col {...titleCol}>
						<Form.Item
							label="Interest reserves" // 储备金
							name="reserves_amount"
							rules={[{ required: true }]}
						>
							<InputNumAmount placeholder="Interest reserves" />
						</Form.Item>
					</Col>
				</>
			)}
			<Col {...titleCol}>
				<Form.Item
					label="As is value" // zeebee中Estimated current value（ltv）
					name="ltv"
					rules={[{ required: true }]}
				>
					<InputNumAmount placeholder="As is value" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item label="After repair value" name="arv" rules={[{ required: true }]}>
					<InputNumAmount placeholder="After repair value" />
				</Form.Item>
			</Col>
			<Col {...titleCol}>
				<Form.Item
					label="Total requested loan amount" // zeebee中Total loan amount
					name="total_amount"
					rules={[{ required: true }]}
				>
					<InputNumAmount placeholder="Total requested loan amount" />
				</Form.Item>
			</Col>
			{isPurchase && (
				<>
					<Col {...layoutCol} md={12}>
						<Form.Item
							label="Assignment fee"
							name="has_assignment_fee"
							rules={[{ required: true }]}
						>
							<Radio.Group options={yesNoOptions} style={{ width: '100%' }} />
						</Form.Item>
					</Col>
					{has_assignment_fee && (
						<Col {...layoutCol} md={12}>
							<Form.Item
								label="Assignment fee"
								name="assignment_fee_amount"
								rules={[
									{
										required: true,
										message: 'Please enter assignment fee amount!',
									},
									{
										validator: (_, value) =>
											validateAssignmentFeeAmount(value, purchase_amount),
									},
								]}
							>
								<InputNumAmount placeholder="Enter assignment fee amount" />
							</Form.Item>
						</Col>
					)}
				</>
			)}
		</Row>
	);
};
export default LoanInfo;
