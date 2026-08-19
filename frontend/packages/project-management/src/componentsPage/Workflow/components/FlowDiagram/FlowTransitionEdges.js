'use client';

import { getCenter, getTransitionKey } from './flowLayout';
import { computeTransitionPath } from './computeTransitionPath';

/**
 * @param {object} props
 */
export default function FlowTransitionEdges({ svgRef, maxX, maxY, flow }) {
	const {
		transitions,
		states,
		positions,
		selectedTransitionKey,
		connecting,
		mousePos,
		selectTransition,
		openConfigDrawer,
	} = flow;

	return (
		<svg
			ref={svgRef}
			style={{
				position: 'absolute',
				top: 0,
				left: 0,
				width: maxX,
				height: maxY,
				pointerEvents: 'none',
			}}
		>
			<defs>
				<marker
					id="arrowhead"
					markerWidth="10"
					markerHeight="7"
					refX="9"
					refY="3.5"
					orient="auto"
				>
					<polygon points="0 0, 10 3.5, 0 7" fill="#999" />
				</marker>
				<marker
					id="arrowhead-selected"
					markerWidth="10"
					markerHeight="7"
					refX="9"
					refY="3.5"
					orient="auto"
				>
					<polygon points="0 0, 10 3.5, 0 7" fill="#1890ff" />
				</marker>
				<marker
					id="arrowhead-active"
					markerWidth="10"
					markerHeight="7"
					refX="9"
					refY="3.5"
					orient="auto"
				>
					<polygon points="0 0, 10 3.5, 0 7" fill="#1890ff" />
				</marker>
			</defs>

			{transitions.map((transition, index) => {
				const pathResult = computeTransitionPath({
					transition,
					transitions,
					states,
					positions,
				});
				if (!pathResult) return null;

				const { pathD, midX, midY } = pathResult;
				const transitionKey = getTransitionKey(transition, index);
				const isSelected = selectedTransitionKey === transitionKey;

				return (
					<g key={transitionKey}>
						<path
							d={pathD}
							fill="none"
							stroke="transparent"
							strokeWidth={14}
							style={{ pointerEvents: 'stroke', cursor: 'pointer' }}
							onClick={(e) => {
								e.stopPropagation();
								selectTransition(transition, transitionKey);
							}}
						/>
						<path
							d={pathD}
							fill="none"
							stroke={isSelected ? '#1890ff' : '#bbb'}
							strokeWidth={isSelected ? 2.5 : 1.5}
							markerEnd={isSelected ? 'url(#arrowhead-selected)' : 'url(#arrowhead)'}
							style={{ pointerEvents: 'none' }}
						/>
						<circle
							cx={midX}
							cy={midY}
							r={9}
							fill="#fff"
							stroke={isSelected ? '#1890ff' : '#aaa'}
							strokeWidth={1.5}
							style={{ pointerEvents: 'all', cursor: 'pointer' }}
							onClick={(e) => {
								e.stopPropagation();
								selectTransition(transition, transitionKey);
								openConfigDrawer(transition, e);
							}}
						/>
						<line
							x1={midX - 4}
							y1={midY + 4}
							x2={midX + 4}
							y2={midY - 4}
							stroke={isSelected ? '#1890ff' : '#aaa'}
							strokeWidth={1}
							style={{ pointerEvents: 'none' }}
						/>
						<line
							x1={midX - 2}
							y1={midY + 4}
							x2={midX + 6}
							y2={midY - 2}
							stroke={isSelected ? '#1890ff' : '#aaa'}
							strokeWidth={1}
							style={{ pointerEvents: 'none' }}
						/>
						<line
							x1={midX - 6}
							y1={midY + 2}
							x2={midX + 2}
							y2={midY - 6}
							stroke={isSelected ? '#1890ff' : '#aaa'}
							strokeWidth={1}
							style={{ pointerEvents: 'none' }}
						/>
					</g>
				);
			})}

			{connecting ? (
				<line
					x1={getCenter(positions, connecting.fromId).x}
					y1={getCenter(positions, connecting.fromId).y}
					x2={mousePos.x}
					y2={mousePos.y}
					stroke="#1890ff"
					strokeWidth={2}
					strokeDasharray="6,3"
					markerEnd="url(#arrowhead-active)"
				/>
			) : null}
		</svg>
	);
}
