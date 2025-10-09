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
	...props
}) => {
	const expandedRowRender = (row, index, b) => {
		console.log(row, b);

		return (
			<List
				itemLayout="horizontal"
				dataSource={[{ title: '1' }]}
				renderItem={(item, index) => (
					<List.Item>
						<List.Item.Meta title={<a href="https://ant.design">{item.title}</a>} />
						<Space>
							<Button>view</Button>
							<Button>delete</Button>
						</Space>
					</List.Item>
				)}
			/>
		);
	};
	return (
		// <Table
		// 	bordered={false}
		// 	dataSource={dataSource}
		// 	loading={loading}
		// 	columns={[...columns, Table.EXPAND_COLUMN]}
		// 	pagination={false}
		// 	expandable={{ expandedRowRender }}
		// 	className={[styles.documentList, props?.className].join(' ')}
		// 	{...props}
		// />
		<div className={[styles.documentList, props?.className].join(' ')}>
			<div
				className="flex py-4 px-4 font-medium gap-4"
				style={{ color: 'rgba(0, 0, 0, 0.65)' }}
			>
				{columns.map((i) => (
					<div
						key={i.key}
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
					return <DocumentItem key={index} item={item} columns={columns}></DocumentItem>;
				})}
			</div>
		</div>
	);
};
export default DocumentList;
