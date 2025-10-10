import React from 'react';
import NewTable from '@/components/NewTable';
import styles from './index.module.scss';
import { Button, List, Space, Table } from 'antd';
import DocumentItem from './Item';
const DocumentList = ({
	dataSource = [
		{
			total: 0,
		},
	],
	loading,
	columns,
	isHide,
	...props
}) => {
	return (
		<div className={[styles.documentList, props?.className].join(' ')}>
			<div
				className="flex py-4 px-4 font-medium gap-4"
				style={{ color: 'rgba(0, 0, 0, 0.65)' }}
			>
				{columns.map((i, index) => (
					<div
						key={index}
						style={{
							width: i.width,
							minWidth: i.minWidth,
							flex: i.minWidth ? 1 : '0 0 auto',
						}}
						className={[i?.ellipsis && styles.ellipsis].join(' ')}
					>
						{i.title}
					</div>
				))}
			</div>
			<div>
				{dataSource.map((item, index) => {
					return (
						<DocumentItem
							key={index}
							item={item}
							columns={columns}
							isHide={isHide}
						></DocumentItem>
					);
				})}
			</div>
		</div>
	);
};
export default DocumentList;
