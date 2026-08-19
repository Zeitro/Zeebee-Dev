'use client';

import { useCallback } from 'react';
import { useTranslations } from 'next-intl';
import {
	getWorkflowStateColorHex,
	WORKFLOW_STATE_TYPE_OPTIONS,
} from '@repo/project-management/enum/workflow';
import { NODE_H, NODE_W, portDirs, portOffsets } from './flowLayout';

/**
 * @param {object} props
 */
export default function FlowStateNodes({ flow }) {
	const t = useTranslations();

	const {
		states,
		positions,
		dragging,
		activeNodeId,
		hoveredNodeId,
		setHoveredNodeId,
		connecting,
		startDrag,
		clickNode,
		endConnection,
		startFromPort,
	} = flow;

	const getStateTypeLabel = useCallback(
		(type) => {
			const option = WORKFLOW_STATE_TYPE_OPTIONS.find((item) => item.value === type);
			return option ? t(option.labelKey) : type;
		},
		[t],
	);

	return states.map((state) => {
		const pos = positions.get(state.id);
		if (!pos) return null;
		const bgColor = getWorkflowStateColorHex(state.color);
		const isActive = activeNodeId === state.id;
		const isHovered = hoveredNodeId === state.id;
		const isConnecting = connecting?.fromId === state.id;
		const showPorts = isActive || isHovered || isConnecting;

		return (
			<div
				key={state.id}
				style={{
					position: 'absolute',
					left: pos.x,
					top: pos.y,
					width: NODE_W,
					height: NODE_H,
					background: '#fff',
					border: `2px solid ${isActive ? bgColor : isHovered ? bgColor : '#ddd'}`,
					borderRadius: 10,
					display: 'flex',
					alignItems: 'center',
					justifyContent: 'center',
					cursor:
						dragging?.nodeId === state.id
							? 'grabbing'
							: connecting
								? 'crosshair'
								: 'grab',
					userSelect: 'none',
					boxShadow: isActive
						? `0 0 0 3px ${bgColor}33, 0 2px 8px rgba(0,0,0,0.1)`
						: '0 1px 4px rgba(0,0,0,0.06)',
					transition: 'border-color 0.2s, box-shadow 0.2s',
					zIndex: dragging?.nodeId === state.id ? 10 : isActive ? 2 : 1,
				}}
				onMouseDown={(e) => {
					if (!connecting) startDrag(state.id, e);
				}}
				onClick={(e) => {
					if (!dragging) clickNode(state.id, e);
				}}
				onMouseUp={() => {
					if (connecting && connecting.fromId !== state.id) {
						endConnection(state.id);
					}
				}}
				onMouseEnter={() => setHoveredNodeId(state.id)}
				onMouseLeave={() => setHoveredNodeId(null)}
			>
				<div
					style={{
						position: 'absolute',
						left: -6,
						top: '50%',
						transform: 'translateY(-50%)',
						width: 12,
						height: 12,
						borderRadius: '50%',
						background: bgColor,
						border: '2px solid #fff',
					}}
				/>
				<div style={{ textAlign: 'center' }}>
					<div style={{ fontWeight: 600, fontSize: 13, color: '#333' }}>{state.name}</div>
					<div style={{ fontSize: 10, color: '#999' }}>
						{state.is_initial
							? t('workflow.editor.stateInitial')
							: state.is_terminal
								? t('workflow.editor.stateTerminal')
								: getStateTypeLabel(state.type)}
					</div>
				</div>

				{showPorts
					? portDirs.map((dir) => {
							const portPos = portOffsets[dir];
							const portX = portPos.dx * NODE_W;
							const portY = portPos.dy * NODE_H;
							return (
								<div
									key={dir}
									style={{
										position: 'absolute',
										left: portX - 6,
										top: portY - 6,
										width: 12,
										height: 12,
										borderRadius: '50%',
										background: '#fff',
										border: `2px solid ${bgColor}`,
										cursor: 'crosshair',
										zIndex: 3,
										transition: 'transform 0.15s',
									}}
									onMouseDown={(e) => startFromPort(state.id, e)}
									onMouseEnter={(e) => {
										e.currentTarget.style.transform = 'scale(1.3)';
									}}
									onMouseLeave={(e) => {
										e.currentTarget.style.transform = 'scale(1)';
									}}
								/>
							);
						})
					: null}
			</div>
		);
	});
}
