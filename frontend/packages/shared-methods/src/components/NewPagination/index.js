import { Pagination } from 'antd';
import React from 'react';
import styles from './index.module.scss';
export const NewPagination = ({ total, current, pageSize, onChange = () => {}, ...props }) => {
	return (
		<Pagination
			showQuickJumper={false}
			showSizeChanger={false}
			total={total}
			onChange={onChange}
			current={current}
			pageSize={pageSize}
			{...props}
			className={[styles.page, 'shrink-0', props.className].join(' ')}
		/>
	);
};
