'use client';

import { Button, Space } from 'antd';
import { CloseOutlined, ReloadOutlined } from '@ant-design/icons';

/**
 * @param {object} props
 * @param {string} props.toolbarHint
 * @param {boolean} props.showCancel
 * @param {() => void} props.onCancel
 * @param {() => void} props.onResetLayout
 * @param {(key: string) => string} props.t
 */
export default function FlowDiagramToolbar({
	toolbarHint,
	showCancel,
	onCancel,
	onResetLayout,
	t,
}) {
	return (
		<div
			style={{
				marginBottom: 12,
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'space-between',
			}}
		>
			<span style={{ fontSize: 13, color: '#666' }}>{toolbarHint}</span>
			<Space>
				{showCancel ? (
					<Button size="small" icon={<CloseOutlined />} onClick={onCancel}>
						{t('common.cancel')}
					</Button>
				) : null}
				<Button size="small" icon={<ReloadOutlined />} onClick={onResetLayout}>
					{t('workflow.editor.resetLayout')}
				</Button>
			</Space>
		</div>
	);
}
