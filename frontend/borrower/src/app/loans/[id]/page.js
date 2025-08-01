'use client';
import { Anchor, Button, Col, Form, Row, Space, Steps, Tabs } from 'antd';
import React, { useEffect, useState } from 'react';
import CreateLoan from '../components/Create/Loan';
import CreateProperty from '../components/Create/property';
import CreateTerm from '../components/Create/term';
// import styles from './loans.module.scss';
import { useRouter } from 'next/navigation';
import { createLoanSteps } from '@/enum/loan';
import styles from './index.module.scss';
const editLoanPage = () => {
	const [currentStep, setCurrentStep] = useState(2);
	const [form] = Form.useForm();
	const stage = Form.useWatch('stage', form);

	// next/Add loan
	const handleNext = async () => {
		const values = await form.validateFields();
		if (currentStep < 2) {
			setCurrentStep((v) => v + 1);
			return;
		}
		// edit loan
	};

	return (
		<Row className={['h-full', styles.editLoan].join(' ')}>
			<Col flex={'200px'} className="h-full border-r border-[#DBDCDE]">
				<Anchor
					getContainer={() => document.getElementById('editLoan-content')}
					replace
					affix={true}
					onClick={() => {}}
					items={createLoanSteps.map((item) => ({ ...item, href: `#${item.key}` }))}
				/>
			</Col>
			<Col
				id="editLoan-content"
				flex="calc(100% - 200px)"
				className="h-full overflow-hidden overflow-y-auto px-8"
			>
				<Form form={form} name="basic" layout="vertical" autoComplete="off" size="large">
					<div id="loanDetail">
						<p className="py-6 font-medium text-lg/7">Loan detail</p>
						<CreateLoan form={form} xs={24} sm={24} lg={12} gutter={24} />
					</div>
					<div id="propertyDetail">
						<p className="py-6 font-medium text-lg/7 border-t border-base-border">
							Property detail
						</p>
						<CreateProperty stage={stage} xs={24} sm={24} lg={12} gutter={24} />
					</div>
					<div id="loanTerms">
						<p className="py-6 font-medium text-lg/7 border-t border-base-border">
							Loan terms
						</p>
						<CreateTerm stage={stage} xs={24} sm={24} lg={12} gutter={24} />
					</div>
				</Form>
				<div className="flex justify-end py-4">
					<Button type="primary" onClick={handleNext}>
						Save change
					</Button>
				</div>
			</Col>
		</Row>
	);
};
export default editLoanPage;
