'use client';

import React from 'react';
import { Space, Spin } from 'antd';
import { NoData } from '@repo/shared-methods';
import CollapsibleDocumentRow from './CollapsibleDocumentRow';
import { getColumnFlexStyle } from './columnHelpers';
import styles from './index.module.scss';

/**
 * 可折叠文档列表：父级（文档类型）与子级（文件）共用列网格，按 rowPlacement 占位对齐。
 *
 * 列配置：
 * - rowPlacement: 'parent' | 'file' | 'both'（默认 parent）
 * - parentDataIndex / fileDataIndex / dataIndex
 * - render(value, record, { level, dataIndex })：仅父级
 * - fileRender(value, defaultNode, file, parentRecord)
 *
 * 子级默认：首列文件名 pill + 末列 View / Delete。
 * 仅单侧展示的列：另一侧用空占位元素对齐，不显示 `--`。
 */
const CollapsibleDocumentList = ({
	columns = [],
	dataSource = [],
	loading,
	className,
	showHeader = true,
	filesKey = 'files',
	rowKey = 'id',
	hideFileDelete = false,
	onDeleteFile,
	onAfterDelete,
	fileKey = 'file_url',
	empty,
}) => {
	const rootClassName = [styles.list, className].filter(Boolean).join(' ');

	return (
		<div className={rootClassName}>
			{showHeader && columns.length > 0 && (
				<div className={styles.columnHeader}>
					{columns.map((column, index) => (
						<div
							key={column.key ?? column.dataIndex ?? index}
							style={getColumnFlexStyle(column)}
							className={[column?.ellipsis && 'min-w-0 truncate']
								.filter(Boolean)
								.join(' ')}
						>
							{column.title}
						</div>
					))}
				</div>
			)}

			{dataSource.length > 0 ? (
				<Space direction="vertical" size={12} className="min-w-full">
					{dataSource.map((item, index) => (
						<CollapsibleDocumentRow
							key={item?.[rowKey] ?? index}
							item={item}
							columns={columns}
							filesKey={filesKey}
							rowKey={rowKey}
							hideFileDelete={hideFileDelete}
							onDeleteFile={onDeleteFile}
							onAfterDelete={onAfterDelete}
							fileKey={fileKey}
						/>
					))}
				</Space>
			) : (
				<div className={styles.empty}>
					{empty ?? <NoData>{loading && <Spin />}</NoData>}
				</div>
			)}
		</div>
	);
};

export default CollapsibleDocumentList;
