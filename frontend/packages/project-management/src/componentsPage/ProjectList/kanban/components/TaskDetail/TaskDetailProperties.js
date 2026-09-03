/**
 * 任务详情属性区：左标签 + 右控件的属性行布局（对齐 Linear / Notion 风格）。
 */

import {
	CalendarOutlined,
	CheckCircleOutlined,
	FieldTimeOutlined,
	FlagOutlined,
	TagOutlined,
	UserOutlined,
} from '@ant-design/icons';
import { CommonFormDate } from '@repo/shared-methods';
import { Form, InputNumber, Select } from 'antd';
import { getWorkflowStateColor } from '../../../../../enum/workflow';
import { AssigneeChips } from './AssigneeChips';

/** 属性行内控件统一底色与圆角 */
const FIELD_SURFACE_CLASS =
	'w-full [&_.ant-select-selector]:!bg-base-text2 [&_.ant-select-selector]:!rounded-md [&_.ant-select-selector]:!border-0 [&_.ant-picker]:!bg-base-text2 [&_.ant-picker]:!rounded-md [&_.ant-picker]:!border-0 [&_.ant-input-number]:!bg-base-text2 [&_.ant-input-number]:!rounded-md [&_.ant-input-number]:!border-0 [&_.ant-input-number-input]:!bg-transparent';

/**
 * 属性行：左侧标签 + 右侧控件
 *
 * 注意：该组件会被放到外层网格（2~3列）中，因此这里使用 `1fr` 让控件自适应列宽，
 * 避免“值区域过宽”的问题。
 *
 * @param {{
 *   label: string,
 *   icon?: import('react').ReactNode,
 *   children: import('react').ReactNode,
 * }} props
 */
function PropertyRow({ label, icon, children }) {
	return (
		<div className="grid grid-cols-[86px_1fr] gap-x-3 gap-y-1 items-center py-1.5">
			<span className="text-[13px] text-base-label inline-flex items-center gap-1.5 shrink-0">
				{icon}
				{label}
			</span>
			<div className="min-w-0">{children}</div>
		</div>
	);
}

/**
 * @param {{
 *   t: (key: string) => string,
 *   workflow?: object|null,
 *   taskTypes?: Array,
 *   priorityOptions?: Array,
 *   memberOptions?: Array,
 *   members?: Array,
 * }} props
 */
export function TaskDetailProperties({
	t,
	workflow,
	taskTypes = [],
	priorityOptions = [],
	projectId,
}) {
	return (
		<div className="rounded-xl border border-base-border bg-white p-2 mb-6">
			<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-x-3 gap-y-0.5">
				<PropertyRow
					label={t('kanban.activityFieldStatus')}
					icon={<CheckCircleOutlined className="text-xs" />}
				>
					<Form.Item name="workflow_state_id" noStyle>
						<Select variant="borderless" size="small" className={FIELD_SURFACE_CLASS}>
							{(workflow?.states || []).map((state) => (
								<Select.Option key={state.id} value={state.id}>
									<span
										className="inline-block w-2 h-2 rounded-full mr-1.5 align-middle"
										style={{
											background: getWorkflowStateColor(state.color),
										}}
									/>
									{state.name}
								</Select.Option>
							))}
						</Select>
					</Form.Item>
				</PropertyRow>

				<PropertyRow label={t('kanban.type')} icon={<TagOutlined className="text-xs" />}>
					<Form.Item name="type" noStyle>
						<Select
							variant="borderless"
							size="small"
							className={FIELD_SURFACE_CLASS}
							allowClear
							placeholder={t('kanban.type')}
						>
							{taskTypes.map((item) => (
								<Select.Option key={item.id} value={item.name}>
									{item.name}
								</Select.Option>
							))}
						</Select>
					</Form.Item>
				</PropertyRow>

				<PropertyRow
					label={t('projectList.priority')}
					icon={<FlagOutlined className="text-xs" />}
				>
					<Form.Item name="priority" noStyle>
						<Select
							variant="borderless"
							size="small"
							className={FIELD_SURFACE_CLASS}
							options={priorityOptions}
						/>
					</Form.Item>
				</PropertyRow>

				<PropertyRow
					label={t('kanban.startDate')}
					icon={<CalendarOutlined className="text-xs" />}
				>
					<Form.Item name="start_date" noStyle>
						<CommonFormDate
							variant="borderless"
							size="small"
							placeholder={t('kanban.startDate')}
							suffixIcon={null}
							className={FIELD_SURFACE_CLASS}
						/>
					</Form.Item>
				</PropertyRow>

				<PropertyRow
					label={t('kanban.dueDate')}
					icon={<CalendarOutlined className="text-xs" />}
				>
					<Form.Item name="due_date" noStyle>
						<CommonFormDate
							variant="borderless"
							size="small"
							placeholder={t('kanban.dueDate')}
							suffixIcon={null}
							className={FIELD_SURFACE_CLASS}
						/>
					</Form.Item>
				</PropertyRow>

				<PropertyRow
					label={t('kanban.estimatedHours')}
					icon={<FieldTimeOutlined className="text-xs" />}
				>
					<Form.Item name="estimated_hours" noStyle>
						<InputNumber
							variant="borderless"
							size="small"
							placeholder={t('kanban.estimatedHours')}
							min={0}
							step={0.5}
							suffix={t('kanban.hoursSuffix')}
							className={FIELD_SURFACE_CLASS}
						/>
					</Form.Item>
				</PropertyRow>

				<div className="sm:col-span-2 lg:col-span-3">
					<PropertyRow
						label={t('kanban.assignee')}
						icon={<UserOutlined className="text-xs" />}
					>
						<Form.Item name="assignee_ids" noStyle>
							<AssigneeChips projectId={projectId} />
						</Form.Item>
					</PropertyRow>
				</div>
			</div>
		</div>
	);
}
