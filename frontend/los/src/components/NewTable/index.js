import { Table } from 'antd';
import React from 'react';
import styles from './index.module.scss';
const NewTable = ({ dataSource, columns, ...props }) => {
	return (
		<Table
			dataSource={dataSource}
			columns={columns}
			rowKey="id"
			{...props}
			className={[styles.newTable, props?.className].join(' ')}
		/>
	);
};
export default NewTable;
