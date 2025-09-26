'use client';
import styles from './index.module.scss';
import { Form, Menu, Segmented, Steps, Tabs } from 'antd';
import React, { useEffect } from 'react';
import useCreate from '../components/CreateZeeBee/hooks/useCreate';
import { assemblyAddress } from '@/lib/format';
import { flattenArray } from '@/lib/common';
import { stageColor, pendingStatusColor, loanDetailsFields } from '@/enum/loan';
import HomeSvg from '@/assets/images/home.svg';
import _ from 'lodash';
const ZeeBeeDetails = () => {
	const {
		scrollRef,
		currentStep,
		setCurrentStep,
		setIsChange,
		lenderDetailsSteps,
		form,
		initialValues,
		rowData,
		setRowData,
		lenderDetailsComponents,
	} = useCreate('details');

	useEffect(() => {
		setCurrentStep(lenderDetailsSteps[0].key);
	}, []);

	const getLenderMenu = (list = lenderDetailsSteps, isChild) => {
		const iconFn = (item, isChild, index) => {
			if (isChild) {
				return (
					<div className="!w-[2px] !min-w-[2px] bg-base-border h-full mr-[10px]"></div>
				);
			}
			return (
				<div
					className="bg-blue-dark w-6 h-6 rounded-full flex !items-center justify-center !text-white font-normal"
					style={{
						paddingTop: item.icon ? 0 : '2px',
					}}
				>
					{item.icon || index + 1}
				</div>
			);
		};
		return list.map((item, index) => {
			let obj = {
				key: item.key,
				label: item.title,
				icon: iconFn(item, isChild, index),
			};
			if (item.children) {
				obj.children = getLenderMenu(item.children, true);
			}
			return obj;
		});
	};

	return (
		<div className="flex flex-col h-full">
			<div className="flex items-center flex-shrink-0 pt-4 pb-5 px-4">
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
			<div className="flex-shrink-0 flex items-center gap-6 px-4 pb-4 text-base/[22px]">
				{loanDetailsFields?.[rowData?.stage]?.map((item, index) => (
					<div className="flex items-center gap-2 " key={index}>
						<div className="text-[#4e5969]">{item.label}:</div>
						<div className="text-[#000]">
							{item.render?.(rowData[item.value]) ?? (rowData?.[item.value] || '--')}
						</div>
					</div>
				))}
			</div>
			<div className="flex-1 overflow-hidden flex items-center bg-gray-light border-t border-base-border">
				<Menu
					className={`flex-shrink-0 !border-none !w-[260px] h-full !bg-transparent !pl-4 !py-3 ${styles.menu}`}
					openKeys={lenderDetailsSteps.map((item) => item.key)}
					selectedKeys={[currentStep]}
					mode="inline"
					theme="light"
					items={getLenderMenu()}
					onSelect={({ key }) => {
						setCurrentStep(key);
						scrollRef.current.scrollTo({
							top: 0,
						});
					}}
					size="large"
				/>
				<div
					className={`flex-1 overflow-hidden overflow-y-auto bg-white h-full`}
					ref={scrollRef}
					// style={{
					// 	paddingTop: stepContent.length === currentStep + 1 ? 16 : 24,
					// }}
				>
					<Form
						form={form}
						name="basic"
						layout="vertical"
						autoComplete="off"
						className=""
						size="large"
						// disabled={true}
						initialValues={initialValues}
						scrollToFirstError
						onValuesChange={() => {
							setIsChange(true);
						}}
						onFieldsChange={(changeFields) => {
							console.log('changeFields', changeFields);

							// 只有输入验证通过了才可更新值
							const isErr = changeFields.every(
								(v) => v?.errors?.length || v?.warnings?.length || !v?.validated,
							);
							if (!isErr) {
								changeFields.map((v) => {
									if (!v?.errors?.length && !v?.warnings?.length) {
										_.set(rowData, v?.name, v?.value);
									}
								});
								setRowData({ ...rowData });
							}
						}}
					>
						<Tabs
							activeKey={currentStep}
							items={flattenArray(lenderDetailsSteps).map((i) => ({
								key: i.key,
								label: i.title,
								children: lenderDetailsComponents[i.key],
							}))}
							renderTabBar={() => <div> </div>}
							className={` ${styles.tabs}`}
						/>
					</Form>
				</div>
			</div>
			{/* <Segmented
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
			</div> */}
		</div>
	);
};
export default ZeeBeeDetails;
