'use client';
import { Anchor, Button, Col, Form, message, Row, Space, Steps, Tabs } from 'antd';
import React, { useEffect, useState } from 'react';
import CreateLoan from '../components/Create/Loan';
import CreateProperty from '../components/Create/property';
import CreateTerm from '../components/Create/term';
// import styles from './loans.module.scss';
import { usePathname, useParams } from 'next/navigation';
import { createLoanSteps } from '@/enum/loan';
import styles from './index.module.scss';
import { createRequest } from '@/lib/fetch';
import dayjs from 'dayjs';
import { formatDate } from '@/lib/date';
const editLoanPage = () => {
	// const router = useRouter();
	const pathname = usePathname();
	const { id } = useParams();
	const [currentStep, setCurrentStep] = useState(2);
	const [form] = Form.useForm();
	const stage = Form.useWatch('stage', form);
	const { request } = createRequest();

	const getDetails = async () => {
		let res = {
			applicants: [
				{ is_primary: true, type: 'Individual' },
				{ is_primary: false, type: 'Individual' },
			],
		};
		try {
			res = await request(`/loan/${id}`);
			if (res.closing_date) {
				res.closing_date = dayjs(res.closing_date);
			}
			if (res.maturity_date) {
				res.maturity_date = dayjs(res.maturity_date);
			}
			if (res.paid_off_date) {
				res.paid_off_date = dayjs(res.paid_off_date);
			}
			if (!res.applicants?.length < 2) {
				res.applicants = [
					{ is_primary: true, type: 'Individual', ...(res.applicants?.[0] || {}) },
					{ is_primary: false, type: 'Individual', ...(res.applicants?.[1] || {}) },
				];
			}
			form.setFieldsValue({
				...res,
			});
		} catch (error) {
			form.setFieldsValue({
				...res,
			});
		}
	};

	// edit loan
	const handleNext = async () => {
		const values = await form.validateFields();
		if (values.closing_date) {
			values.closing_date = formatDate(values.closing_date);
		}
		if (values.maturity_date) {
			values.maturity_date = formatDate(values.maturity_date);
		}
		if (values.paid_off_date) {
			values.paid_off_date = formatDate(values.paid_off_date);
		}
		await request(`/loan/${id}`, {
			method: 'PUT',
			body: values,
		});
		message.open({
			type: 'success',
			content: 'Update successful！',
		});
	};

	useEffect(() => {
		getDetails();
	}, []);

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
						<CreateProperty
							stage={stage}
							xs={24}
							sm={24}
							lg={12}
							gutter={24}
							form={form}
						/>
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
