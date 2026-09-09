'use client';

import { Button, Popconfirm, Space, Tag } from 'antd';
import { DeleteOutlined } from '@ant-design/icons';

/**
 * @param {object} props
 */
export default function FlowTransitionPanel({
	selectedTransition,
	states,
	canEdit,
	onOpenConfig,
	onDelete,
	t,
}) {
	if (!selectedTransition) return null;

	const fromName = states.find((s) => s.id === selectedTransition.from_state_id)?.name;
	const toName = states.find((s) => s.id === selectedTransition.to_state_id)?.name;

	return (
		<div
			style={{
				marginTop: 12,
				padding: '12px 16px',
				background: '#e6f7ff',
				borderRadius: 8,
				border: '1px solid #91d5ff',
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'space-between',
			}}
		>
			<div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
				<div
					style={{
						width: 8,
						height: 8,
						borderRadius: '50%',
						background: '#1890ff',
						flexShrink: 0,
					}}
				/>
				<strong style={{ fontSize: 13 }}>
					{fromName}
					{' → '}
					{toName}
				</strong>
				{selectedTransition.name ? <Tag color="blue">{selectedTransition.name}</Tag> : null}
				<Tag>{selectedTransition.trigger_type}</Tag>
			</div>
			<Space>
				<Button size="small" onClick={onOpenConfig}>
					{t('workflow.editor.setConditions')}
				</Button>
				{canEdit ? (
					<Popconfirm
						title={t('workflow.editor.deleteTransitionConfirm')}
						onConfirm={onDelete}
					>
						<Button danger size="small" icon={<DeleteOutlined />}>
							{t('common.delete')}
						</Button>
					</Popconfirm>
				) : null}
			</Space>
		</div>
	);
}
