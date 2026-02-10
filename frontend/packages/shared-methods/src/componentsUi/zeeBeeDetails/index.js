'use client';
import styles from './index.module.scss';
import { Form, Menu, Skeleton, Tabs } from 'antd';
import React, { useEffect, useRef, useState } from 'react';
import useLenderDetails from './components/useLenderDetails';
import { assemblyAddress, changeUnitByPropertyType } from '@repo/shared-methods/utils';
import { stageColor, pendingStatusColor, loanDetailsFields } from '@repo/shared-methods/enum/loan';
import HomeSvg from '@repo/shared-methods/assets/images/home.svg';
import _ from 'lodash';
import { CheckOutlined } from '@ant-design/icons';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
// loan详情页面组件（lender端/admin端使用）
const ZeeBeeDetails = ({ isAllReadonly }) => {
	const {
		scrollRef,
		currentStep,
		setCurrentStep,
		setIsChange,
		lenderDetailsSteps,
		form,
		initialValues,
		rowData,
		newRowData,
		setNewRowData,
		disabled,
		tabList,
		loanProgress, // 当前loan根据stage和pendingStatus计算的进度
		unitInfo,
		setUnitInfo,
		detailsLoading,
	} = useLenderDetails({ isAllReadonly });
	let [lenderMenu, setLenderMenu] = useState([]);
	// 计算总步骤，方便跟loanProgress比较。
	let count = useRef({ index: -1, key: '' });

	useEffect(() => {
		setCurrentStep(lenderDetailsSteps[0].key);
	}, [lenderDetailsSteps, setCurrentStep]);

	const getLenderMenu = (list = lenderDetailsSteps, isChild) => {
		const iconFn = (item, isChild, index) => {
			if (isChild) {
				return (
					<div
						className={[
							'!w-[2px] !min-w-[2px] h-full mr-[10px]',
							loanProgress.index > count.current.index
								? 'bg-[#081265]'
								: 'bg-base-border',
						].join(' ')}
					></div>
				);
			}
			const isChildNotFullyCompleted = item?.children?.some(
				(item) => item.key === loanProgress.selectTabKey,
			);

			const isFinish =
				!item.icon && loanProgress.index > count.current.index && !isChildNotFullyCompleted;
			return (
				<div
					className={[
						'w-6 h-6 rounded-full flex !items-center justify-center  font-normal',
						isFinish ? 'bg-[#E0E7FF]' : 'bg-[#081265]',
						isFinish ? '!text-[#081265]' : '!text-[#E0E7FF]',
					].join(' ')}
					style={{
						paddingTop: item.icon ? 0 : '2px',
					}}
				>
					{item.icon || (isFinish ? <CheckOutlined /> : index)}
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
			} else {
				count.current.index = count.current.index + 1;
				count.current.key = item.key;
			}
			return obj;
		});
	};

	useDeferredSyncEffect({
		deferred: () => {
			const arr = getLenderMenu();
			count.current.index = -1;
			setLenderMenu(arr);
		},
		deps: [lenderDetailsSteps, loanProgress.index, loanProgress.selectTabKey],
		isUseMicrotask: true,
	});

	const skeletonProps = {
		title: false,
		active: true,
		loading: detailsLoading,
	};
	return (
		<div className="flex flex-col flex-1 overflow-hidden">
			<Skeleton
				{...skeletonProps}
				className="pt-4 pb-[23px] px-4"
				paragraph={{ rows: 1, width: '30%' }}
				avatar={{ active: true, size: 16 }}
			>
				<div className="flex items-center flex-shrink-0 pt-4 pb-5 px-4">
					<HomeSvg />
					<p className="text-xl/[24px] font-bold pl-3 pr-5">
						{assemblyAddress(rowData?.loan_property?.address, true)}
					</p>
					<div
						className="px-2 text-white font-medium rounded-sm mr-2 text-sm/[22px]"
						style={{ background: stageColor[rowData?.stage] }}
					>
						{rowData?.stage === 'Draft' ? 'Pre approval' : rowData?.stage}
					</div>
					{rowData?.stage === 'Pending' && (
						<div
							style={{
								background: pendingStatusColor?.[rowData?.pending_status]?.bg,
								color: pendingStatusColor?.[rowData?.pending_status]?.color,
							}}
							className="text-sm/[22px] px-2 rounded-sm"
						>
							{rowData?.pending_status}
						</div>
					)}
				</div>
			</Skeleton>
			<Skeleton
				{...skeletonProps}
				className="px-4 pb-[22px]"
				paragraph={{ rows: 1, width: '60%' }}
			>
				<div className="flex-shrink-0 flex items-center gap-6 px-4 pb-4 text-base/[22px]">
					{loanDetailsFields?.[rowData?.stage]?.map((item, index) => (
						<div className="flex items-center gap-2 " key={index}>
							<div className="text-[#4e5969]">{item.label}:</div>
							<div className="text-[#000]">
								{item.render?.(rowData[item.value]) ??
									(rowData?.[item.value] || '--')}
							</div>
						</div>
					))}
				</div>
			</Skeleton>

			<div className="flex-1 overflow-hidden flex bg-gray-light border-t border-base-border">
				<Skeleton
					{...skeletonProps}
					className="pl-4 pr-2 py-6 !w-[260px] flex-shrink-0"
					paragraph={{ rows: 10 }}
					title={true}
				>
					<Menu
						className={`flex-shrink-0 !border-none !w-[260px] h-full !bg-transparent !pl-4 !py-3 ${styles.menu}`}
						openKeys={lenderDetailsSteps.map((item) => item.key)}
						selectedKeys={[currentStep]}
						mode="inline"
						theme="light"
						items={lenderMenu}
						onSelect={({ key }) => {
							setCurrentStep(key);
							scrollRef.current.scrollTo({
								top: 0,
							});
						}}
						size="large"
						expandIcon={() => null}
					/>
				</Skeleton>
				<Skeleton
					paragraph={{ rows: 12 }}
					className="px-6 py-6 bg-background"
					{...skeletonProps}
					title={true}
				>
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
							disabled={isAllReadonly || disabled}
							initialValues={initialValues}
							scrollToFirstError
							onValuesChange={(v) => {
								setIsChange(true);
								// 根据property_type设置unit
								const property_type = v?.loan_property?.future_property_type;
								if (property_type) {
									const info = changeUnitByPropertyType(property_type);
									const unit = { ...unitInfo };
									unit.future_unit = info;
									setUnitInfo({ ...unit });
									form.setFieldValue(['loan_property', 'future_unit'], info.unit);
								}
							}}
							onFieldsChange={(changeFields) => {
								console.log('changeFields', changeFields);

								// 只有输入验证通过了才可更新值
								const isErr = changeFields.every(
									(v) =>
										v?.errors?.length || v?.warnings?.length || !v?.validated,
								);
								if (!isErr) {
									changeFields.map((v) => {
										if (!v?.errors?.length && !v?.warnings?.length) {
											_.set(newRowData, v?.name, v?.value);
										}
									});
									setNewRowData({ ...newRowData });
								}
							}}
						>
							<Tabs
								activeKey={currentStep}
								items={tabList}
								renderTabBar={() => <div> </div>}
								className={` ${styles.tabs}`}
								destroyOnHidden
							/>
						</Form>
					</div>
				</Skeleton>
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
