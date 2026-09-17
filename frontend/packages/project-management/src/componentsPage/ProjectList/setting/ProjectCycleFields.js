/**
 * 项目设置表单里的周期 / 开始 / 截止。改周期或开始日时，把截止日期重算为开始 + 周期。
 */
'use client';

import { DatePicker, Form, InputNumber, Select, Space } from 'antd';
import { useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { dayjs, formatStr } from '@repo/shared-methods/utils';
import { PROJECT_CYCLE_UNITS, addProjectCycle } from './cycle';

/**
 * @param {{ required?: boolean, stacked?: boolean }} [props]
 * 创建项目时周期 / 开始日必填，并在弹框内纵向排开；设置页保持可选、三列。
 */
export function ProjectCycleFields({ required = false, stacked = false } = {}) {
	const t = useTranslations();
	const form = Form.useFormInstance();
	const cycleAmount = Number(Form.useWatch('cycle_months', form)) || 3;
	const cycleUnit = Form.useWatch('cycle_unit', form) || 'month';
	const startDate = Form.useWatch('start_date', form);

	const unitOptions = useMemo(
		() => PROJECT_CYCLE_UNITS.map((item) => ({ value: item.value, label: t(item.labelKey) })),
		[t],
	);

	const applyDefaultEndDate = () => {
		const start = form.getFieldValue('start_date');
		const amount = Number(form.getFieldValue('cycle_months')) || 3;
		const unit = form.getFieldValue('cycle_unit') || 'month';
		if (!start) return;
		form.setFieldValue('end_date', addProjectCycle(dayjs(start), amount, unit));
	};

	return (
		<div className={stacked ? 'flex flex-col' : 'grid grid-cols-3 gap-4'}>
			<Form.Item label={t('projectSetting.cycle')} required={required}>
				<Space.Compact block>
					<Form.Item
						name="cycle_months"
						noStyle
						rules={
							required
								? [{ required: true, message: t('projectSetting.cycleRequired') }]
								: undefined
						}
					>
						<InputNumber
							min={1}
							precision={0}
							step={1}
							className="!w-full"
							onChange={applyDefaultEndDate}
						/>
					</Form.Item>
					<Form.Item name="cycle_unit" noStyle>
						<Select
							options={unitOptions}
							className="!w-[88px]"
							onChange={applyDefaultEndDate}
						/>
					</Form.Item>
				</Space.Compact>
			</Form.Item>
			<Form.Item
				name="start_date"
				label={t('projectSetting.startDate')}
				rules={
					required
						? [{ required: true, message: t('projectSetting.startDateRequired') }]
						: undefined
				}
			>
				<DatePicker className="w-full" format={formatStr} onChange={applyDefaultEndDate} />
			</Form.Item>
			<Form.Item name="end_date" label={t('projectSetting.endDate')}>
				<DatePicker
					className="w-full"
					format={formatStr}
					disabledDate={(current) => {
						if (!current) return false;
						const start = startDate ? dayjs(startDate) : null;
						if (!start) return false;
						const max = addProjectCycle(start, cycleAmount, cycleUnit);
						return (
							current.isBefore(start.startOf('day')) ||
							current.isAfter(max.endOf('day'))
						);
					}}
				/>
			</Form.Item>
		</div>
	);
}
