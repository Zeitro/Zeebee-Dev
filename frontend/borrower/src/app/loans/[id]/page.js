'use client';
import { Anchor, Button, Col, Form, message, Row, Skeleton } from 'antd';
import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import CreateLoan from '../components/Create/Loan';
import CreateProperty from '../components/Create/property';
import CreateTerm from '../components/Create/term';
import { useParams, useRouter } from 'next/navigation';
import { createLoanSteps } from '@/enum/loan';
import styles from './index.module.scss';
import {
	assemblyAddress,
	changeUnitByPropertyType,
	createRequest,
	eventBus,
} from '@repo/shared-methods/utils';
import { useGlobalStore } from '@/store/global';
import QuoteBtn from '../components/Quote';
import { PermissionSlot } from '@repo/shared-methods';
import { PERMISSIONS } from '@config/permissions';

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
	const [detailsLoading, setDetailsLoading] = useState(true);
	const { changeSidebarCollapsed } = useGlobalStore();
	const property_type = Form.useWatch(['loan_property', 'property_type'], form);
	const [unitInfo, setUnitInfo] = useState({});
	let isInit = useRef(true);

	const getDetails = async () => {
		let res = initialValues;
		try {
			setDetailsLoading(true);
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
			setDetailsLoading(false);
		} catch {
			form.setFieldsValue({
				...res,
			});
			setDetailsLoading(false);
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

	const skeletonProps = {
		active: true,
		loading: detailsLoading,
		paragraph: { rows: 4 },
		className: 'pt-6',
	};

	return (
		<Row className={['h-full flex-1 overflow-hidden', styles.editLoan].join(' ')}>
			<Col flex={'200px'} className="h-full border-r border-[#DBDCDE]">
				<Skeleton
					title={false}
					className="pt-6 px-6"
					paragraph={{ rows: 3 }}
					active
					loading={detailsLoading}
				>
					<Anchor
						getContainer={() => document.getElementById('editLoan-content')}
						replace
						affix={true}
						onClick={() => {}}
						items={createLoanSteps.map((item) => ({ ...item, href: `#${item.key}` }))}
					/>
				</Skeleton>
			</Col>
			<Col
				id="editLoan-content"
				flex="calc(100% - 200px)"
				className="h-full overflow-hidden overflow-y-auto px-8"
			>
				<Form
					form={form}
					name="Edit-3rdParty-basic"
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
					<Skeleton {...skeletonProps}>
						<div id="loanDetail">
							<p className="py-6 font-medium text-lg/7">Loan detail</p>
							<CreateLoan form={form} xs={24} sm={24} lg={12} gutter={24} />
						</div>
					</Skeleton>
					<Skeleton {...skeletonProps}>
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
					</Skeleton>
					<Skeleton {...skeletonProps}>
						<div id="loanTerms">
							<p className="py-6 font-medium text-lg/7 border-t border-base-border">
								Loan terms
							</p>
							<CreateTerm
								stage={stage}
								xs={24}
								sm={24}
								lg={12}
								gutter={24}
								form={form}
							/>
						</div>
					</Skeleton>
				</Form>
				<Skeleton
					{...skeletonProps}
					paragraph={{ rows: 1, width: '100px' }}
					title={{ width: '100px' }}
				>
					<div className="flex justify-end py-4 gap-4">
						<PermissionSlot code={PERMISSIONS?.LOAN?.QUOTE}>
							<QuoteBtn
								disabled={!(stage === 'Pending' && pending_status === 'Quoting')}
								getRowData={async () => await form.validateFields()}
								isGetOrganization={true}
							/>
						</PermissionSlot>
						<Button type="primary" onClick={handleNext} loading={loading}>
							Save change
						</Button>
					</div>
				</Skeleton>
			</Col>
		</Row>
	);
};

export default EditLoanPage;
