import React from 'react';
import styles from './index.module.scss';
import { Space } from 'antd';
import DocumentItem from './Item';
import { NoData } from '@repo/shared-methods';
const DocumentList = ({
	dataSource = [
		{
			total: 0,
		},
	],
	columns,
	isHide,
	handleClickOpt,
	deleteLoading,
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
			<Space direction="vertical" className="w-full">
				{dataSource?.length > 0 ? (
					dataSource.map((item, index) => {
						return (
							<DocumentItem
								key={index}
								item={item}
								columns={columns}
								isHide={isHide}
								handleClickOpt={handleClickOpt}
								deleteLoading={deleteLoading}
							></DocumentItem>
						);
					})
				) : (
					<NoData className="bg-[#f2f2f2]" />
				)}
			</Space>
		</div>
	);
};
export default DocumentList;
