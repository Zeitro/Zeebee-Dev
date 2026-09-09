'use client';

import { Button, Input, Space, Segmented } from 'antd';
import {
	ArrowLeftOutlined,
	PlusOutlined,
	ApartmentOutlined,
	UnorderedListOutlined,
} from '@ant-design/icons';
import { useTranslations } from 'next-intl';

/**
 * 编辑器顶栏：返回、名称编辑、视图切换、添加状态。
 */
export function WorkflowEditorHeader({
	workflowName,
	editingName,
	nameInputValue,
	nameInputRef,
	viewMode,
	canEdit,
	onBack,
	onStartEditName,
	onNameChange,
	onSaveName,
	onViewModeChange,
	onAddState,
}) {
	const t = useTranslations();

	return (
		<div className="flex items-center justify-between mb-4 flex-shrink-0">
			<Space>
				<Button icon={<ArrowLeftOutlined />} onClick={onBack}>
					{t('workflow.editor.back')}
				</Button>
				<h2 className="m-0 flex items-center gap-2 text-lg font-semibold">
					{editingName ? (
						<>
							<Input
								ref={nameInputRef}
								value={nameInputValue}
								onChange={(e) => onNameChange(e.target.value)}
								onPressEnter={onSaveName}
								className="!w-[200px]"
								autoFocus
							/>
							<Button type="primary" size="small" onClick={onSaveName}>
								{t('workflow.editor.save')}
							</Button>
						</>
					) : (
						<span
							className="cursor-pointer border-b border-dashed border-gray-border"
							onClick={onStartEditName}
						>
							{workflowName}
						</span>
					)}
				</h2>
			</Space>
			<Space>
				<Segmented
					value={viewMode}
					onChange={onViewModeChange}
					options={[
						{
							label: t('workflow.editor.viewFlow'),
							value: 'flow',
							icon: <ApartmentOutlined />,
						},
						{
							label: t('workflow.editor.viewList'),
							value: 'list',
							icon: <UnorderedListOutlined />,
						},
					]}
				/>
				<Button
					type="primary"
					icon={<PlusOutlined />}
					onClick={onAddState}
					disabled={!canEdit}
				>
					{t('workflow.editor.addState')}
				</Button>
			</Space>
		</div>
	);
}
