'use client';
import { Anchor, Button, Col, Form, message, Row } from 'antd';
import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import CreateLoan from '../components/Create/Loan';
import CreateProperty from '../components/Create/property';
import CreateTerm from '../components/Create/term';
import { useParams, useRouter } from 'next/navigation';
import { createLoanSteps } from '@/enum/loan';
import styles from './index.module.scss';
import {
	eventBus,
	assemblyAddress,
	changeUnitByPropertyType,
	createRequest,
} from '@repo/shared-methods/utils';
import { useGlobalStore } from '@/store/global';
import QuoteBtn from '../components/Quote';
const initialValues = {
	applicants: [{ is_primary: true, type: 'Individual' }],
};
const EditLoanPage = () => {
	const { id } = useParams();
	const router = useRouter();
	const [form] = Form.useForm();
	const stage = Form.useWatch('stage', form);
	const pending_status = Form.useWatch(['pending_status'], form);
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const { changeSidebarCollapsed } = useGlobalStore();
	const property_type = Form.useWatch(['loan_property', 'property_type'], form);
	const [unitInfo, setUnitInfo] = useState({});
	let isInit = useRef(true);

	const getDetails = async () => {
		let res = initialValues;
		try {
			res = await request(`/loan/${id}`);
			eventBus.emit('changeCrumbs', {
				2: `Edit Loan(${assemblyAddress(res?.loan_property?.address)})`,
			});
			if (!res.applicants?.length) {
				res.applicants = initialValues?.applicants;
			}
			form.setFieldsValue({
				...res,
			});
		} catch {
			form.setFieldsValue({
				...res,
			});
		}
	};

	// edit loan
	const handleNext = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			values.applicants = values?.applicants.filter(
				(i) => i.first_name || i.mid_name || i.last_name,
			);
			await request(`/loan/${id}`, {
				method: 'PUT',
				body: values,
			});
			setLoading(false);
			message.open({
				type: 'success',
				content: 'Update successful！',
			});
			router.back();
		} catch {
			setLoading(false);
		}
	};

	useEffect(() => {
		getDetails();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	useLayoutEffect(() => {
		changeSidebarCollapsed(true);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	// 通过property_type改变unit
	useEffect(() => {
		if (property_type && isInit.current) {
			const info = changeUnitByPropertyType(property_type);
			setUnitInfo(info);
			isInit.current = false;
		}
	}, [property_type]);

	return (
		<Row className={['h-full flex-1 overflow-hidden', styles.editLoan].join(' ')}>
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
				<Form
					form={form}
					name="basic"
					disabled={loading}
					layout="vertical"
					autoComplete="off"
					size="large"
					initialValues={initialValues}
					// onValuesChange={(v) => {
					// 	// 根据property_type设置unit
					// 	const property_type = v?.loan_property?.property_type;
					// 	if (property_type) {
					// 		const info = changeUnitByPropertyType(property_type);
					// 		form.setFieldValue(
					// 			['loan_property', 'property_base', 'unit'],
					// 			info.unit,
					// 		);
					// 		setUnitInfo({ ...info });
					// 	}
					// }}
				>
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
							unitInfo={unitInfo}
						/>
					</div>
					<div id="loanTerms">
						<p className="py-6 font-medium text-lg/7 border-t border-base-border">
							Loan terms
						</p>
						<CreateTerm stage={stage} xs={24} sm={24} lg={12} gutter={24} form={form} />
					</div>
				</Form>
				<div className="flex justify-end py-4 gap-4">
					<QuoteBtn
						disabled={!(stage === 'Pending' && pending_status === 'Quoting')}
						onSave={handleNext}
						getRowData={async () => await form.validateFields()}
					/>
					<Button type="primary" onClick={handleNext} loading={loading}>
						Save change
					</Button>
				</div>
			</Col>
		</Row>
	);
};

export default EditLoanPage;
