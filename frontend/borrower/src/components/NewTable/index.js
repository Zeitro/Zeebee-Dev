'use client';
import { Table } from 'antd';
import React, { useLayoutEffect, useMemo, useRef, useState } from 'react';
import styles from './index.module.scss';
// import { debounce } from 'lodash';
// const limitHeight = 200;
const NewTable = ({ dataSource, columns, pagination, rowKey = 'id', scroll, ...props }) => {
	const [selectedRowKey, setSelectedRowKey] = useState(null);

	// const [scrollY, setScrollY] = useState(limitHeight);

	// useLayoutEffect(() => {
	// 	const handleResize = debounce(() => {
	// 		const tableDom = document?.querySelector(`.${styles.newTable}`);
	// 		const headerDom = document?.querySelector(`.${styles.newTable} .ant-table-header`);
	// 		if (!tableDom) return;
	// 		// 1. 获取容器高度
	// 		const containerHeight = tableDom.clientHeight;
	// 		// 2. 计算表头高度（假设表头高度为56px，分页高度为64px，可根据实际调整）
	// 		const headerHeight = headerDom?.offsetHeight || 55;
	// 		const paginationHeight = 64;
	// 		// 3. 可用高度
	// 		const maxTableBodyHeight = containerHeight - headerHeight - paginationHeight; // 16为padding等
	// 		// 4. 计算实际内容高度
	// 		const rowHeight = 55; // 你的每行高度，按实际调整
	// 		const bodyHeight = dataSource?.length
	// 			? dataSource?.length * rowHeight + 16
	// 			: limitHeight;
	// 		// 5. 取最小值
	// 		console.log(dataSource,bodyHeight,maxTableBodyHeight,'--',containerHeight);

	// 		setScrollY(Math.min(maxTableBodyHeight, bodyHeight));
	// 	}, 300);
	// 	handleResize();
	// 	window.addEventListener('resize', handleResize);

	// 	return () => {
	// 		window.removeEventListener('resize', handleResize);
	// 	};
	// }, [dataSource]);
	return (
		<Table
			dataSource={dataSource}
			columns={columns}
			rowKey={rowKey}
			pagination={pagination}
			{...props}
			scroll={{ x: 'max-content', y: 'calc(100% - 55px)', ...(scroll || {}) }}
			className={[styles.newTable, props?.className].join(' ')}
			onRow={(record) => {
				const moreConfig = props?.onRow?.(record);
				return {
					...moreConfig,
					onClick: (e) => {
						setSelectedRowKey(record[rowKey]);
						props?.onRow?.(record)?.onClick?.(e);
					}, // 点击行
				};
			}}
			rowClassName={(record) => (record[rowKey] === selectedRowKey ? styles.selectedRow : '')}
		/>
	);
};
export default NewTable;
