'use client';
import { Form, Steps, Tabs } from 'antd';
import React from 'react';
import styles from './index.module.scss';
import useCreate from '../components/CreateZeeBee/hooks/useCreate';
import { changeUnitByPropertyType } from '@repo/shared-methods/utils';
const CreateZeeBeeLoan = () => {
	const {
		getTabsItem,
		scrollRef,
		stepContent,
		currentStep,
		stepNum,
		stepBtn,
		setIsChange,
		createZeeBeeSteps,
		form,
		initialValues,
		loading,
		setUnitInfo,
		unitInfo,
	} = useCreate();

	return (
		<div className="flex flex-col h-full flex-1 overflow-hidden">
			<div className="px-12 py-4 shrink-0 flex items-center justify-center border-b border-base-border">
				<Steps size="default" current={currentStep} items={createZeeBeeSteps} />
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
						disabled={loading}
						initialValues={initialValues}
						scrollToFirstError
						onValuesChange={(v) => {
							setIsChange(true);
							// 根据property_type设置unit
							const property_type =
								v?.loan_property?.property_type ||
								v?.loan_property?.future_property_type;
							if (property_type) {
								const info = changeUnitByPropertyType(property_type);
								const fieldKey = v?.loan_property?.property_type
									? 'unit'
									: 'future_unit';
								const moreFieldKey =
									fieldKey === 'unit' ? ['property_base', fieldKey] : [fieldKey];
								const unit = { ...unitInfo };
								unit[fieldKey] = info;
								setUnitInfo({ ...unit });
								form.setFieldValue(['loan_property', ...moreFieldKey], info.unit);
							}
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
