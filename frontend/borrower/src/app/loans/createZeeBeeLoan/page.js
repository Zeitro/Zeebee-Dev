'use client';
import { Button, Col, Form, message, Space, Steps, Tabs } from 'antd';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import styles from './index.module.scss';
import { useRouter, useSearchParams } from 'next/navigation';
import { createZeeBeeSteps } from '@/enum/loan';
import { createRequest } from '@/lib/fetch';
import ApplicantInfo from '../components/CreateZeeBee/ApplicantInfo';
import BorrowerInfo from '../components/CreateZeeBee/BorrowerInfo';
import PropertyInfo from '../components/CreateZeeBee/PropertyInfo';
import LoanDetail from '../components/CreateZeeBee/Loan';
import Experience from '../components/CreateZeeBee/Experience';
import ThirdPartyService from '../components/CreateZeeBee/ThirdPartyService';
import Declaration from '../components/CreateZeeBee/Declaration';
import Demographics from '../components/CreateZeeBee/Demographics';
import ReviewAndSign from '../components/CreateZeeBee/ReviewAndSign';
import _ from 'lodash';
import { flattenObjectPaths } from '@/lib/common';
const CreateZeeBeeLoan = () => {
	const router = useRouter();
	const searchParams = useSearchParams();
	const api = `/property/${searchParams.get('propertyId')}`;
	const [form] = Form.useForm();
	const scrollRef = useRef(null);
	const [currentStep, setCurrentStep] = useState(0); //顶部-总步骤
	const [stepNum, setStepNum] = useState(0); //侧边步骤
	const [loading, setLoading] = useState(false);
	const [errorList, setErrorList] = useState([]);
	const { request } = createRequest();

	// 从property详情跳转的创建loan
	const getProperty = async () => {
		const data = await request(api);
		const newData = {
			loan_property: {
				// property_type: data.property_type,
				property_base: {
					unit: data.unit,
					address: data.address,
				},
			},
		};
		form.setFieldsValue(newData);
	};

	// 获取property
	useEffect(() => {
		if (searchParams.has('propertyId')) getProperty();
	}, []);

	const stepContent = [
		[
			<ApplicantInfo key="0-ApplicantInfo" form={form} />,
			<BorrowerInfo key="1-BorrowerInfo" form={form} />,
		],
		[
			<PropertyInfo
				key="0-PropertyInfo"
				form={form}
				isRequestProperty={searchParams.has('propertyId')}
			/>,
			<LoanDetail key="1-LoanDetail" form={form} />,
		],
		[
			<Experience key="0-Experience" form={form} />,
			<ThirdPartyService key="1-ThirdPartyService" form={form} />,
		],
		[
			<Declaration key="0-Declaration" form={form} />,
			<Demographics key="1-Demographics" form={form} />,
		],
		<ReviewAndSign key="0-ReviewAndSign" form={form} />,
	];

	// 将stepContent的多维数组，转换成一维，key对应的就是当前步骤（currentStep + stepNum）
	const getTabsItem = () => {
		let newArr = [];
		stepContent.map((item, index) => {
			if (Array.isArray(item)) {
				item.map((item2, index2) => {
					newArr.push({
						key: `${index}-${index2}`,
						label: '',
						children: item2,
					});
				});
			} else {
				newArr.push({
					key: `${index}-0`,
					label: '',
					children: item,
				});
			}
		});
		return newArr;
	};

	// 获取下一步步骤类型（Submit创建按钮/NextMaxStep下一大步骤/NextMinStep下一小步骤）
	const diffSteps = () => {
		const lastStep = stepContent[stepContent.length - 1];
		// 某大步骤下有多个小步骤（数组），如为数组需判断是否最后一小步骤
		const isLastStepNumByArr = lastStep?.length ? stepNum + 1 === lastStep?.length : true;
		if (currentStep + 1 === stepContent.length && isLastStepNumByArr) {
			return 'Submit';
		}
		const parentStepData = createZeeBeeSteps[currentStep];
		const childArr = parentStepData?.children;
		if (childArr ? stepNum + 1 === childArr?.length : !childArr) {
			return 'NextMaxStep';
		}

		return 'NextMinStep';
	};

	// // next/Add loan
	const handleNext = async () => {
		try {
			const a = form.getFieldsValue();
			let fieldsKeys = true;
			// 解决：没填数据就触发了错误验证，又直接切换到上一步骤，又点击下一步时，validateFields会报错，提示刚才错误验证的字段没有填，导致下一页切换不了
			if (errorList.length) {
				errorList.map((i) => {
					_.unset(a, i); // 拿到错误的字段key，删除这项内容
				});
				fieldsKeys = flattenObjectPaths(a); // 将多层嵌套对象转换为二维路径数组
			}
			const values = await form.validateFields(fieldsKeys);
			setErrorList([]);
			const isType = diffSteps();
			if (isType === 'Submit') {
				setLoading(true);
				const body = _.cloneDeep(values);

				body.loan_borrower.tin_type = values?.loan_borrower?.tin?.tin_type;
				body.loan_borrower.tin = values?.loan_borrower?.tin?.tin;
				await request('/loan/apply', {
					method: 'POST',
					body,
				});
				message.open({
					type: 'success',
					content: 'Success in creation！',
				});
				setLoading(false);
				// 		router.back();
				return;
			}
			if (isType === 'NextMaxStep') {
				setCurrentStep((v) => v + 1);
				setStepNum(0);
			} else {
				setStepNum((v) => v + 1);
			}

			scrollRef.current.scrollTo({
				top: 0,
			});
		} catch (error) {
			if (error?.errorFields) {
				setErrorList(error?.errorFields.map((item) => item.name));
			}
			setLoading(false);
		}
	};

	// 下一步按钮文案
	const getNextBtnText = () => {
		const isType = diffSteps();
		const parentStepData = createZeeBeeSteps[currentStep];
		switch (isType) {
			case 'Submit':
				return 'Submit application';
			case 'NextMaxStep':
				return `Next ${createZeeBeeSteps[currentStep + 1]?.title}`;
			case 'NextMinStep':
				return `Next ${parentStepData?.children[stepNum + 1]?.title}`;
			default:
				break;
		}
	};

	const handleBack = () => {
		if (errorList.length) {
			// 光设置重新验证不行，需要结合validateFields处的操作才可避免
			form.setFields(
				errorList.map((i) => ({
					name: i,
					errors: [],
					touched: false,
					validated: false,
				})),
			);
			form.resetFields(errorList);
		}
		if (stepNum) {
			setStepNum((v) => (v - 1 < 0 ? 0 : v - 1));
		} else {
			const prevStepData = stepContent[currentStep - 1];
			setCurrentStep((v) => (v - 1 < 0 ? 0 : v - 1));
			setStepNum(prevStepData?.length - 1 < 0 ? 0 : prevStepData?.length - 1);
		}
	};

	const stepBtn = useMemo(() => {
		const cancelBtn = <Button onClick={() => router.back()}>Cancel</Button>;
		const backBtn = !(!currentStep && !stepNum) && (
			<Button className="w-[90px]" onClick={handleBack}>
				Back
			</Button>
		);
		const nextBtn = (
			<Button type="primary" onClick={handleNext} loading={loading}>
				{getNextBtnText()}
			</Button>
		);
		if (currentStep === stepContent.length - 1) {
			return (
				<div className="w-[198px] m-auto flex flex-col gap-4">
					{nextBtn}
					<Space size={24}>
						{cancelBtn}
						{backBtn}
					</Space>
				</div>
			);
		}
		return (
			<Col xs={24} sm={24} className="py-6 !flex justify-between items-center">
				{cancelBtn}
				<Space size={8}>
					{backBtn}
					{nextBtn}
				</Space>
			</Col>
		);
	}, [currentStep, stepNum, loading, errorList]);
	// createZeeBeeSteps.map((i, index) => ({
	// 							key: index,
	// 							label: i.title,
	// 							children: (
	// 								<>
	// 									{stepContent[index]?.length
	// 										? stepContent[index]?.[stepNum]
	// 										: stepContent[index]}
	// 									{stepBtn}
	// 								</>
	// 							),
	// 						}))

	return (
		<div className="flex flex-col h-full">
			<div className="px-12 py-4 shrink-0 flex items-center justify-center border-b border-base-border">
				<Steps
					size="default"
					current={currentStep}
					items={createZeeBeeSteps}
					// className="max-w-[880px]"
				/>
			</div>

			<div className="flex-1 flex items-stretch overflow-hidden">
				{createZeeBeeSteps[currentStep]?.children && (
					<div className="w-[250px] py-6 px-6 shrink-0 border-r border-base-border">
						<Steps
							size="small"
							progressDot
							direction="vertical"
							current={stepNum}
							items={createZeeBeeSteps[currentStep]?.children}
						/>
					</div>
				)}
				<div
					className={`flex-1 overflow-hidden overflow-y-auto pb-8 px-8`}
					ref={scrollRef}
					style={{
						paddingTop: stepContent.length === currentStep + 1 ? 24 : 32,
					}}
				>
					<Form
						form={form}
						name="basic"
						layout="vertical"
						autoComplete="off"
						className=""
						size="large"
						// disabled={loading}
						initialValues={{
							applicants: [{ type: 'Individual', is_primary: true }],
						}}
					>
						<Tabs
							activeKey={`${currentStep}-${stepNum}`}
							items={getTabsItem()}
							renderTabBar={() => <div> </div>}
							className={` ${styles.tabs}`}
						/>
					</Form>
					{stepBtn}
				</div>
			</div>
		</div>
	);
};
export default CreateZeeBeeLoan;
