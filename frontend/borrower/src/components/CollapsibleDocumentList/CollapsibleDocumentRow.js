'use client';

import React, { useState } from 'react';
import { Collapse } from 'antd';
import { NoData } from '@repo/shared-methods';
import { displayValue } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import {
	ROW_LEVEL,
	getColumnCellClassName,
	getColumnFlexStyle,
	getRecordValue,
	isPlaceholderCell,
	renderPlaceholder,
} from './columnHelpers';
import { DefaultFileActions, DefaultFileNameCell, DefaultParentNameCell } from './DefaultFileCells';
import styles from './index.module.scss';

function renderParentCell(column, record, columnIndex, totalColumns, open) {
	const value = getRecordValue(record, column, ROW_LEVEL.PARENT);
	const isFirstColumn = columnIndex === 0;
	const defaultNameCell = isFirstColumn ? (
		<DefaultParentNameCell
			record={record}
			fieldKey={column.parentDataIndex ?? column.dataIndex}
			open={open}
		/>
	) : null;

	if (column.render) {
		return column.render(value, record, {
			level: ROW_LEVEL.PARENT,
			dataIndex: column.parentDataIndex ?? column.dataIndex,
		});
	}

	return defaultNameCell ?? displayValue(value, column.formatter);
}

function renderFileCell(column, file, parentRecord, columnIndex, totalColumns, actionProps) {
	if (isPlaceholderCell(column, ROW_LEVEL.FILE)) {
		return renderPlaceholder(ROW_LEVEL.FILE);
	}

	const value = getRecordValue(file, column, ROW_LEVEL.FILE);
	const isFirstColumn = columnIndex === 0;
	const defaultNameCell = isFirstColumn ? (
		<DefaultFileNameCell file={file} fileKey={actionProps.fileKey} />
	) : null;
	const isActionColumn = columnIndex === totalColumns - 1;
	const defaultActionCell = isActionColumn ? (
		<DefaultFileActions file={file} {...actionProps} />
	) : null;
	const defaultCell =
		defaultNameCell ?? defaultActionCell ?? displayValue(value, column.formatter);

	if (column.fileRender) {
		return column.fileRender(value, defaultCell, file, parentRecord);
	}

	// render 仅用于父级；子级无 fileRender 时走 defaultCell（末列 DefaultFileActions 等）
	return defaultCell;
}

function RowCells({ columns, level, record, parentRecord, actionProps, open }) {
	return (
		<div className="flex w-full items-center gap-4">
			{columns.map((column, columnIndex) => {
				const cellClassName = getColumnCellClassName(
					column,
					level,
					columnIndex,
					columns.length,
				);

				let content;
				if (isPlaceholderCell(column, level)) {
					content = renderPlaceholder(level);
				} else if (level === ROW_LEVEL.PARENT) {
					content = renderParentCell(column, record, columnIndex, columns.length, open);
				} else {
					content = renderFileCell(
						column,
						record,
						parentRecord,
						columnIndex,
						columns.length,
						actionProps,
					);
				}

				return (
					<div
						key={column.key ?? column.dataIndex ?? columnIndex}
						style={getColumnFlexStyle(column, level === ROW_LEVEL.PARENT)}
						className={[cellClassName, 'min-w-0'].filter(Boolean).join(' ')}
					>
						{content}
					</div>
				);
			})}
		</div>
	);
}

const CollapsibleDocumentRow = ({
	item,
	columns,
	filesKey = 'files',
	rowKey = 'id',
	hideFileDelete,
	onDeleteFile,
	onAfterDelete,
	fileKey = 'file_url',
}) => {
	const [activeKey, setActiveKey] = useState([]);
	const files = item?.[filesKey] ?? [];
	const panelKey = item?.[rowKey];

	useDeferredSyncEffect({
		deferred: () => {
			setActiveKey(files.length ? [panelKey] : []);
		},
		deps: [files.length, panelKey],
		delay: 0,
	});

	const actionProps = {
		hideDelete: hideFileDelete,
		onDeleteFile,
		onAfterDelete,
		fileKey,
	};

	return (
		<Collapse
			bordered={false}
			collapsible="header"
			className={styles.collapse}
			activeKey={activeKey}
			onChange={(e) => {
				setActiveKey(e?.map((item) => Number(item)));
			}}
			items={[
				{
					key: panelKey,
					showArrow: false,
					classNames: {
						header: styles.collapseHeader,
					},
					label: (
						<RowCells
							columns={columns}
							level={ROW_LEVEL.PARENT}
							record={item}
							open={activeKey.includes(panelKey)}
						/>
					),
					children: files.length ? (
						<div className={styles.fileList}>
							{files.map((file, index) => (
								<div key={file.id ?? file.key ?? index} className={styles.fileRow}>
									<RowCells
										columns={columns}
										level={ROW_LEVEL.FILE}
										record={file}
										parentRecord={item}
										actionProps={actionProps}
									/>
								</div>
							))}
						</div>
					) : (
						<div className={styles.fileEmpty}>
							<NoData />
						</div>
					),
				},
			]}
		/>
	);
};

export default CollapsibleDocumentRow;
