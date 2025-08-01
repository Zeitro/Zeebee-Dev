import { Pagination } from 'antd';
import React from 'react';
import styles from './index.module.scss';
const NewPagination = ({ total, current, pageSize, onChange=()=>{}, ...props }) => {
	return (
		<Pagination
			showQuickJumper
			showSizeChanger={false}
			total={total}
			onChange={onChange}
			current={current}
			pageSize={pageSize}
			{...props}
			className={[styles.page,'shrink-0', props.className].join(' ')}
		/>
	);
};
export default NewPagination;
