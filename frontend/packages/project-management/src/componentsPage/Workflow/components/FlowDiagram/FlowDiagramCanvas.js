'use client';

import { Button } from 'antd';
import FlowTransitionEdges from './FlowTransitionEdges';
import FlowStateNodes from './FlowStateNodes';

/**
 * @param {object} props
 */
export default function FlowDiagramCanvas({
	containerRef,
	svgRef,
	maxX,
	maxY,
	canvasOffsetY = 0,
	zoom,
	pan,
	isPanning,
	dragging,
	connecting,
	onMouseMove,
	onMouseUp,
	onMouseDown,
	onWheel,
	onClick,
	onZoomOut,
	onZoomIn,
	onZoomReset,
	t,
	flow,
}) {
	const cursor = isPanning
		? 'grabbing'
		: dragging
			? 'grabbing'
			: connecting
				? 'crosshair'
				: 'grab';

	return (
		<div
			ref={containerRef}
			onMouseMove={onMouseMove}
			onMouseUp={onMouseUp}
			onMouseLeave={onMouseUp}
			onMouseDown={onMouseDown}
			onWheel={onWheel}
			onClick={onClick}
			style={{
				position: 'relative',
				overflow: 'hidden',
				background: '#fafafa',
				border: '1px solid #eee',
				borderRadius: 8,
				minHeight: maxY,
				cursor,
			}}
		>
			<div
				style={{
					position: 'absolute',
					bottom: 12,
					right: 12,
					zIndex: 10,
					background: 'rgba(255,255,255,0.9)',
					borderRadius: 6,
					padding: '4px 12px',
					fontSize: 12,
					color: '#666',
					border: '1px solid #eee',
					display: 'flex',
					alignItems: 'center',
					gap: 8,
				}}
			>
				<Button size="small" type="text" onClick={onZoomOut}>
					−
				</Button>
				<span style={{ minWidth: 36, textAlign: 'center' }}>{Math.round(zoom * 100)}%</span>
				<Button size="small" type="text" onClick={onZoomIn}>
					+
				</Button>
				<Button size="small" type="text" onClick={onZoomReset}>
					{t('common.reset')}
				</Button>
			</div>

			<div
				data-flow-canvas="true"
				style={{
					position: 'absolute',
					top: canvasOffsetY,
					left: 0,
					width: maxX,
					height: maxY,
					transformOrigin: '0 0',
					transform: `translate(${pan.x}px, ${pan.y}px) scale(${zoom})`,
				}}
			>
				<FlowTransitionEdges svgRef={svgRef} maxX={maxX} maxY={maxY} flow={flow} />
				<FlowStateNodes flow={flow} />
			</div>
		</div>
	);
}
