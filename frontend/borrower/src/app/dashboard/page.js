'use client';
import { Radio, Select, Space } from 'antd';
import React, { useCallback, useState } from 'react';
import MaturityAlerts from './components/MaturityAlerts';
import LoanStatusCard from './components/LoanStatusCard';
import DifferentStatusLoans from './components/loans';
import { FieldTips } from '@repo/shared-methods';
import { filterTimeCount, filterTimeCountTips } from '@repo/shared-methods/enum';
// 授权gmail，以用户邮箱发送邮件<中转页,接收参数>
const Dashboard = () => {
	const [type, setType] = useState('all');
	const [dayCount, setDayCount] = useState('All');
	const handleChange = (e) => {
		setType(e.target.value);
	};
	return (
		<div>
			<div className="flex items-center justify-between">
				<span>Hi firstName,welcome back!</span>
				<Space>
					<Radio.Group onChange={handleChange} value={type}>
						<Radio.Button value="all">All loans</Radio.Button>
						<Radio.Button value="zeebee">Zeebee loans</Radio.Button>
						<Radio.Button value="3rd">3rd party loans</Radio.Button>
					</Radio.Group>
					<Select
						className="min-w-[100px]"
						onChange={(v) => setDayCount(v)}
						value={dayCount}
						options={filterTimeCount}
						suffixIcon={
							<FieldTips title={filterTimeCountTips} isShowDescription={true} />
						}
					/>
				</Space>
			</div>
			<div className="flex items-stretch flex-wrap">
				<LoanStatusCard dayCount={dayCount} />
				<MaturityAlerts dayCount={dayCount} />
			</div>
			<DifferentStatusLoans stage={'Pending'} dayCount={dayCount} />
		</div>
	);
};
export default Dashboard;
