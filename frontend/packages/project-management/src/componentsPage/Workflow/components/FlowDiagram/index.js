/**
 * 工作流 SVG 流程图：拖拽节点、端口连线、流转配置抽屉。
 * `workflowApi` 由父级注入；布局通过「重置布局」或父级 `layoutKey` remount 刷新。
 */
'use client';

import { useFlowDiagram } from './hooks/useFlowDiagram';
import { DEFAULT_ZOOM } from './hooks/useFlowDiagramInteraction';
import FlowDiagramToolbar from './FlowDiagramToolbar';
import FlowDiagramCanvas from './FlowDiagramCanvas';
import FlowTransitionPanel from './FlowTransitionPanel';
import FlowAddTransitionModal from './FlowAddTransitionModal';
import FlowConfigDrawer from './FlowConfigDrawer';

/**
 * @param {object} props
 * @param {object} props.workflow
 * @param {object} props.workflowApi
 * @param {() => void} props.onUpdate
 * @param {boolean} props.canEdit
 * @param {number} props.userId
 * @param {(state: object|null) => void} [props.onSelectState]
 */
function FlowDiagramInner({
	workflow,
	workflowApi,
	messageApi,
	onUpdate,
	canEdit,
	userId,
	onSelectState,
}) {
	const flow = useFlowDiagram({
		workflow,
		workflowApi,
		messageApi,
		onUpdate,
		canEdit,
		userId,
		onSelectState,
	});

	return (
		<div>
			<FlowDiagramToolbar
				toolbarHint={flow.toolbarHint}
				showCancel={Boolean(flow.connecting || flow.activeNodeId)}
				onCancel={flow.cancelConnection}
				onResetLayout={flow.handleResetLayout}
				t={flow.t}
			/>

			<FlowDiagramCanvas
				containerRef={flow.containerRef}
				svgRef={flow.svgRef}
				maxX={flow.maxX}
				maxY={flow.maxY}
				zoom={flow.zoom}
				pan={flow.pan}
				isPanning={flow.isPanning}
				dragging={flow.dragging}
				connecting={flow.connecting}
				onMouseMove={flow.handleMouseMove}
				onMouseUp={flow.handleMouseUp}
				onMouseDown={flow.handleCanvasMouseDown}
				onWheel={flow.handleCanvasWheel}
				onClick={flow.handleCanvasClick}
				onZoomOut={() => flow.setZoom((z) => Math.max(0.2, z - 0.1))}
				onZoomIn={() => flow.setZoom((z) => Math.min(2, z + 0.1))}
				onZoomReset={() => {
					flow.setZoom(DEFAULT_ZOOM);
					flow.setPan({ x: 0, y: 0 });
				}}
				t={flow.t}
				flow={flow}
			/>

			<FlowTransitionPanel
				selectedTransition={flow.selectedTransition}
				states={flow.states}
				canEdit={flow.canEdit}
				onOpenConfig={() => {
					flow.setConfigTransition(flow.selectedTransition);
					flow.setConfigDrawerOpen(true);
				}}
				onDelete={flow.handleDeleteTransition}
				t={flow.t}
			/>

			<FlowAddTransitionModal
				open={flow.addModalOpen}
				onCancel={() => flow.setAddModalOpen(false)}
				onOk={flow.handleAddTransition}
				newTransition={flow.newTransition}
				states={flow.states}
				transName={flow.transName}
				onTransNameChange={flow.setTransName}
				transTrigger={flow.transTrigger}
				onTransTriggerChange={flow.setTransTrigger}
				t={flow.t}
			/>

			<FlowConfigDrawer
				open={flow.configDrawerOpen}
				onClose={() => flow.setConfigDrawerOpen(false)}
				onSave={flow.handleSaveConfig}
				configTransition={flow.configTransition}
				states={flow.states}
				configData={flow.configData}
				onConfigDataChange={flow.setConfigData}
			/>
		</div>
	);
}

/**
 * @param {object} props
 * @param {object} props.workflow
 * @param {object} props.workflowApi
 * @param {() => void} props.onUpdate
 * @param {boolean} props.canEdit
 * @param {number} props.userId
 * @param {(state: object|null) => void} [props.onSelectState]
 * @param {string|number} [props.layoutKey] 父级传入以在状态变更时 remount 并重算布局
 */
export default function FlowDiagram({ layoutKey, ...props }) {
	return <FlowDiagramInner key={layoutKey} {...props} />;
}
