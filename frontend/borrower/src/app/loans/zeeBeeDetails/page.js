'use client';
import styles from '../createZeeBeeLoan/index.module.scss';
import { Badge, Form, Radio, Segmented, Steps, Tabs } from 'antd';
import React, { use, useEffect, useMemo } from 'react';
import useCreate from '../components/CreateZeeBee/hooks/useCreate';
import { assemblyAddress } from '@/lib/format';
import { stageColor, pendingStatusColor, loanDetailsFields } from '@/enum/loan';
import HomeSvg from '@/assets/images/home.svg';
const ZeeBeeDetails = () => {
	const {
		getTabsItem,
		scrollRef,
		stepContent,
		currentStep,
		setCurrentStep,
		stepNum,
		setIsChange,
		createZeeBeeSteps,
		zeeBeeDetailsTabs,
		form,
		initialValues,
		setStepNum,
		getZeeBeeDetailsTabsItem,
		rowData,
	} = useCreate('details');

	const segmentedList = [
		...zeeBeeDetailsTabs,
		...getTabsItem(createZeeBeeSteps, true).map((i) => ({
			label: i.label,
			value: i.key.replace('-0', ''),
		})),
	];

	useEffect(() => {
		setCurrentStep(zeeBeeDetailsTabs[0].value);
	}, []);

	return (
		<div className="flex flex-col h-full">
			<div className="flex items-center flex-shrink-0 pt-6 pb-5 px-6">
				<HomeSvg />
				<p className="text-xl/[24px] font-bold pl-3 pr-5">
					{assemblyAddress(rowData?.loan_property?.property_base?.address, true)}
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
							{item.render?.(rowData[item.value]) ?? (rowData?.[item.value] || '--')}
						</div>
					</div>
				))}
			</div>
			<Radio.Group
				size="large"
				value={currentStep}
				onChange={({ target: { value } }) => {
					setCurrentStep(value);
					setStepNum(0);
				}}
				className={`flex-shrink-0 ${styles.radioList}`}
			>
				{segmentedList.map((item, index) => (
					<Radio.Button value={item.value} key={index}>
						{item.totalKey ? (
							<Badge
								count={item?.totalKey ? rowData?.[item?.totalKey] || 0 : 0}
								overflowCount={99}
								size="small"
								offset={[3, -3]}
							>
								{item.label}
							</Badge>
						) : (
							item.label
						)}
					</Radio.Button>
				))}
			</Radio.Group>
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
						disabled={!['documents', 'general'].includes(currentStep)}
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
		</div>
	);
};
export default ZeeBeeDetails;
