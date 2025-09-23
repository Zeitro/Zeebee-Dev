'use client';
import styles from '../createZeeBeeLoan/index.module.scss';
import { Button, Col, Form, message, Segmented, Space, Steps, Tabs } from 'antd';
import React, { useEffect } from 'react';
import useCreate from '../components/CreateZeeBee/hooks/useCreate';
import Image from 'next/image';
import { assemblyAddress } from '@/lib/format';
import { stageColor, pendingStatusColor, loanDetailsFields } from '@/enum/loan';
import HomeSvg from '@/assets/images/home.svg';
import General from '../components/CreateZeeBee/General';
const ZeeBeeDetails = () => {
	const {
		getTabsItem,
		scrollRef,
		stepContent,
		currentStep,
		setCurrentStep,
		stepNum,
		stepBtn,
		setIsChange,
		createZeeBeeSteps,
		zeeBeeDetailsTabs,
		form,
		initialValues,
		loading,
		setStepNum,
		getZeeBeeDetailsTabsItem,
		rowData,
	} = useCreate('details');

	useEffect(() => {
		setCurrentStep(zeeBeeDetailsTabs[0].value);
	}, []);

	return (
		<div className="flex flex-col h-full">
			<div className="flex items-center flex-shrink-0 pt-6 pb-5 px-6">
				<HomeSvg />
				<p className="text-xl/[24px] font-bold pl-3 pr-5">
					{assemblyAddress(rowData?.loan_property?.property_base?.address)}
				</p>
				<div
					className="px-2 text-white font-medium rounded-sm mr-2 text-sm"
					style={{ background: stageColor[rowData?.stage] }}
				>
					{rowData?.stage}
				</div>
				{rowData?.stage === 'Pending' && (
					<div
						style={{
							background: pendingStatusColor?.[rowData?.pending_status]?.bg,
							color: pendingStatusColor?.[rowData?.pending_status]?.color,
						}}
						className="text-sm/[24px] px-4 rounded-sm"
					>
						{rowData?.pending_status}
					</div>
				)}
			</div>
			<div className="flex-shrink-0 flex items-center gap-6 px-6 pb-6 text-base/[22px]">
				{loanDetailsFields?.[rowData?.stage]?.map((item, index) => (
					<div className="flex items-center gap-2 " key={index}>
						<div className="text-[#4e5969]">{item.label}:</div>
						<div className="text-[#000]">
							{item.render?.(rowData[item.value]) ?? rowData?.[item.value]}
						</div>
					</div>
				))}
			</div>
			<Segmented
				value={currentStep}
				size="large"
				options={[
					...zeeBeeDetailsTabs,
					...getTabsItem(createZeeBeeSteps, true).map((i) => ({
						label: i.label,
						value: i.key.replace('-0', ''),
					})),
				]}
				onChange={(value) => {
					setCurrentStep(value);
					setStepNum(0);
					console.log(value); // string
				}}
				className={`flex-shrink-0 ${styles.segmented}`}
			/>
			<div className="flex-1 flex items-stretch overflow-hidden">
				{createZeeBeeSteps[currentStep]?.children && (
					<div className="w-[250px] py-6 px-6 shrink-0 border-r border-base-border">
						<Steps
							size="small"
							progressDot
							direction="vertical"
							current={stepNum}
							items={createZeeBeeSteps[currentStep]?.children}
							onChange={(e) => {
								setStepNum(e);
								console.log(e);
							}}
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
						disabled={true}
						initialValues={initialValues}
						scrollToFirstError
						onValuesChange={() => {
							setIsChange(true);
						}}
					>
						<Tabs
							activeKey={`${currentStep}-${stepNum}`}
							items={[...getZeeBeeDetailsTabsItem(), ...getTabsItem()]}
							renderTabBar={() => <div> </div>}
							className={` ${styles.tabs}`}
						/>
					</Form>
				</div>
			</div>
			{/* <Tabs
				// activeKey={`${currentStep}-${stepNum}`}
				items={[...zeeBeeDetailsTabs, ...getTabsItem(createZeeBeeSteps, true, getContent)]}
				className={`flex-1 overflow-hidden ${styles.tabs}`}
				onChange={(e) => {
					console.log(e);
					if (e.includes('-')) {
						setCurrentStep(e.split('-')[0]);
					} else {
						setCurrentStep(e);
					}
				}}
			/> */}
			{/* <div className="px-12 py-4 shrink-0 flex items-center justify-center border-b border-base-border">
				<Steps size="default" current={currentStep} items={createZeeBeeSteps} />
			</div> */}
		</div>
	);
};
export default ZeeBeeDetails;
