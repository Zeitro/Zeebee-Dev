/**
 * 文档解析状态条：已上传 → 解析中 → 解析完成。
 * 纯 UI 展示；status 由父组件传入，parsing 时图标持续旋转。
 */

import React from 'react';
import { Steps } from 'antd';
import {
	CheckCircleOutlined,
	CloudUploadOutlined,
	LoadingOutlined,
	SyncOutlined,
} from '@ant-design/icons';
import styles from './index.module.scss';

/** @typedef {'idle' | 'uploaded' | 'parsing' | 'complete'} ParseStatus */

/**
 * 将业务 status 映射为 antd Steps 的 current + 各步 status。
 * @param {ParseStatus} status
 * @returns {{ current: number, items: object[] }}
 */
function buildStepsConfig(status) {
	const isUploadedDone = status === 'uploaded' || status === 'parsing' || status === 'complete';
	const isParsing = status === 'parsing';
	const isComplete = status === 'complete';

	const current =
		status === 'idle' ? -1 : status === 'uploaded' ? 0 : status === 'parsing' ? 1 : 2;

	return {
		current: Math.max(current, 0),
		items: [
			{
				title: 'Uploaded',
				status: isUploadedDone ? 'finish' : 'wait',
				icon: <CloudUploadOutlined />,
			},
			{
				title: 'Parsing',
				status: isComplete ? 'finish' : isParsing ? 'process' : 'wait',
				// 仅当前步骤为 parsing 时旋转；其余状态用静态 Sync 图标
				icon: isParsing ? <LoadingOutlined spin /> : <SyncOutlined />,
			},
			{
				title: 'Parse complete',
				status: isComplete ? 'finish' : 'wait',
				icon: <CheckCircleOutlined />,
			},
		],
	};
}

/**
 * @param {{ status?: ParseStatus }} props
 */
const ParseStatusBar = ({ status = 'idle' }) => {
	const { current, items } = buildStepsConfig(status);

	return (
		<div className={styles.parseStatusBar}>
			<Steps size="small" current={current} items={items} className={styles.parseSteps} />
		</div>
	);
};

export default ParseStatusBar;
