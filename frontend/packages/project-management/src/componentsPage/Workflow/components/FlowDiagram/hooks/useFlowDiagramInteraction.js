/**
 * 节点拖拽、端口连线与画布平移缩放交互。
 */

import { useState, useRef, useCallback, useEffect } from 'react';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';

/** 画布初始 / 重置缩放比例 */
export const DEFAULT_ZOOM = 0.7;

/**
 * @param {object} params
 * @param {import('antd/es/message/interface').MessageInstance} params.messageApi
 * @param {boolean} params.canEdit
 * @param {(key: string) => string} params.t
 * @param {object[]} params.states
 * @param {Map<number, { x: number; y: number }>} params.positions
 * @param {React.MutableRefObject<HTMLDivElement|null>} params.containerRef
 * @param {number} [params.canvasOffsetY]
 * @param {(state: object|null) => void} [params.onSelectState]
 * @param {number|string|null} [params.selectedStateId]
 * @param {(fromId: number, toId: number) => boolean} params.hasTransition
 * @param {React.Dispatch<React.SetStateAction<object|null>>} params.setSelectedTransition
 * @param {React.Dispatch<React.SetStateAction<string|null>>} params.setSelectedTransitionKey
 * @param {React.Dispatch<React.SetStateAction<boolean>>} params.setAddModalOpen
 * @param {React.Dispatch<React.SetStateAction<{ from: number; to: number }|null>>} params.setNewTransition
 * @param {React.Dispatch<React.SetStateAction<string>>} params.setTransName
 * @param {React.Dispatch<React.SetStateAction<string>>} params.setTransTrigger
 * @param {React.Dispatch<React.SetStateAction<Map<number, { x: number; y: number }>>>} params.setNodePositions
 */
export function useFlowDiagramInteraction({
	canEdit,
	t,
	messageApi,
	states,
	positions,
	containerRef,
	canvasOffsetY = 0,
	onSelectState,
	selectedStateId,
	hasTransition,
	setSelectedTransition,
	setSelectedTransitionKey,
	setAddModalOpen,
	setNewTransition,
	setTransName,
	setTransTrigger,
	setNodePositions,
}) {
	const [dragging, setDragging] = useState(null);
	const [activeNodeId, setActiveNodeId] = useState(selectedStateId ?? null);
	const [hoveredNodeId, setHoveredNodeId] = useState(null);
	const [connecting, setConnecting] = useState(null);
	const [mousePos, setMousePos] = useState({ x: 0, y: 0 });

	const [zoom, setZoom] = useState(DEFAULT_ZOOM);
	const [pan, setPan] = useState({ x: 0, y: 0 });
	const [isPanning, setIsPanning] = useState(false);
	const panStart = useRef({ x: 0, y: 0, panX: 0, panY: 0 });
	const zoomRef = useRef(zoom);
	const panRef = useRef(pan);

	useEffect(() => {
		zoomRef.current = zoom;
		panRef.current = pan;
	}, [pan, zoom]);

	useDeferredSyncEffect({
		isUseMicrotask: true,
		deferred: () => setActiveNodeId(selectedStateId ?? null),
		deps: [selectedStateId],
	});

	const handleMouseMove = useCallback(
		(e) => {
			if (!containerRef.current) return;
			const rect = containerRef.current.getBoundingClientRect();
			const z = zoomRef.current;
			const p = panRef.current;
			const x = (e.clientX - rect.left - p.x) / z;
			const y = (e.clientY - rect.top - p.y) / z - canvasOffsetY;

			if (isPanning) {
				const dx = e.clientX - panStart.current.x;
				const dy = e.clientY - panStart.current.y;
				setPan({ x: panStart.current.panX + dx, y: panStart.current.panY + dy });
				return;
			}

			if (dragging) {
				setNodePositions((prev) => {
					const next = new Map(prev);
					next.set(dragging.nodeId, {
						x: Math.max(0, x - dragging.offsetX),
						y: Math.max(0, y - dragging.offsetY),
					});
					return next;
				});
				return;
			}

			if (connecting) {
				setMousePos({ x, y });
			}
		},
		[canvasOffsetY, connecting, containerRef, dragging, isPanning, setNodePositions],
	);

	const handleMouseUp = useCallback(() => {
		if (dragging) {
			setDragging(null);
		}
		if (isPanning) {
			setIsPanning(false);
		}
	}, [dragging, isPanning]);

	const endConnection = useCallback(
		(toId) => {
			if (!canEdit) return;
			if (!connecting || connecting.fromId === toId) {
				setConnecting(null);
				return;
			}
			if (hasTransition(connecting.fromId, toId)) {
				messageApi.warning(t('workflow.editor.transitionExists'));
				setConnecting(null);
				return;
			}
			setNewTransition({ from: connecting.fromId, to: toId });
			setTransName('');
			setTransTrigger('manual');
			setAddModalOpen(true);
			setConnecting(null);
		},
		[
			canEdit,
			connecting,
			hasTransition,
			messageApi,
			setAddModalOpen,
			setNewTransition,
			setTransName,
			setTransTrigger,
			t,
		],
	);

	const startDrag = useCallback(
		(stateId, e) => {
			if (connecting) return;
			e.stopPropagation();
			e.preventDefault();
			const pos = positions.get(stateId);
			if (!pos || !containerRef.current) return;
			const rect = containerRef.current.getBoundingClientRect();
			const z = zoomRef.current;
			const p = panRef.current;
			const mouseX = (e.clientX - rect.left - p.x) / z;
			const mouseY = (e.clientY - rect.top - p.y) / z - canvasOffsetY;
			setDragging({
				nodeId: stateId,
				offsetX: mouseX - pos.x,
				offsetY: mouseY - pos.y,
			});
		},
		[canvasOffsetY, connecting, containerRef, positions],
	);

	const clickNode = useCallback(
		(stateId, e) => {
			e.stopPropagation();
			if (connecting) {
				endConnection(stateId);
			} else {
				setActiveNodeId(stateId);
				setSelectedTransition(null);
				setSelectedTransitionKey(null);
				onSelectState?.(states.find((s) => String(s.id) === String(stateId)) || null);
			}
		},
		[
			connecting,
			endConnection,
			onSelectState,
			setSelectedTransition,
			setSelectedTransitionKey,
			states,
		],
	);

	const startFromPort = useCallback(
		(stateId, e) => {
			if (!canEdit) return;
			e.stopPropagation();
			e.preventDefault();
			setConnecting({ fromId: stateId });
			if (containerRef.current) {
				const rect = containerRef.current.getBoundingClientRect();
				setMousePos({
					x: e.clientX - rect.left + containerRef.current.scrollLeft,
					y: e.clientY - rect.top + containerRef.current.scrollTop,
				});
			}
		},
		[canEdit, containerRef],
	);

	const cancelConnection = useCallback(() => {
		setConnecting(null);
		setActiveNodeId(null);
		setSelectedTransition(null);
		setSelectedTransitionKey(null);
	}, [setSelectedTransition, setSelectedTransitionKey]);

	const handleCanvasMouseDown = useCallback(
		(e) => {
			const target = e.target;
			if (target === containerRef.current || target.closest('[data-flow-canvas]')) {
				setIsPanning(true);
				panStart.current = { x: e.clientX, y: e.clientY, panX: pan.x, panY: pan.y };
			}
		},
		[containerRef, pan.x, pan.y],
	);

	const handleCanvasWheel = useCallback((e) => {
		e.preventDefault();
		const delta = e.deltaY > 0 ? -0.05 : 0.05;
		setZoom((z) => Math.max(0.2, Math.min(2, z + delta)));
	}, []);

	const handleCanvasClick = useCallback(
		(e) => {
			if (dragging || isPanning) return;
			const el = e.target instanceof Element ? e.target : e.target?.parentElement;
			if (el?.closest?.('button, .ant-modal-wrap, [data-flow-node]')) return;
			if (connecting) setConnecting(null);
			else {
				setSelectedTransition(null);
				onSelectState?.(null);
			}
			setSelectedTransitionKey(null);
		},
		[
			connecting,
			dragging,
			isPanning,
			onSelectState,
			setSelectedTransition,
			setSelectedTransitionKey,
		],
	);

	const selectTransition = useCallback(
		(transition, transitionKey) => {
			setSelectedTransition(transition);
			setSelectedTransitionKey(transitionKey);
			setActiveNodeId(null);
		},
		[setSelectedTransition, setSelectedTransitionKey],
	);

	const toolbarHint = connecting
		? t('workflow.editor.hintConnecting')
		: activeNodeId
			? t('workflow.editor.hintActiveNode')
			: t('workflow.editor.hintDefault');

	return {
		dragging,
		activeNodeId,
		setActiveNodeId,
		hoveredNodeId,
		setHoveredNodeId,
		connecting,
		mousePos,
		zoom,
		setZoom,
		pan,
		setPan,
		isPanning,
		handleMouseMove,
		handleMouseUp,
		startDrag,
		clickNode,
		startFromPort,
		cancelConnection,
		endConnection,
		handleCanvasMouseDown,
		handleCanvasWheel,
		handleCanvasClick,
		selectTransition,
		toolbarHint,
	};
}
