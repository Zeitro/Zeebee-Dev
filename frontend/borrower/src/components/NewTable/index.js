import { Table } from 'antd';
import React from 'react';
import styles from './index.module.scss';
const NewTable = ({ dataSource, columns,pagination,rowKey='id', ...props }) => {
	return (
		<Table
			dataSource={dataSource}
			columns={columns}
			rowKey={rowKey}
            pagination={pagination}
			{...props}
			className={[styles.newTable, props?.className].join(' ')}
		/>
	);
};
export default NewTable;
