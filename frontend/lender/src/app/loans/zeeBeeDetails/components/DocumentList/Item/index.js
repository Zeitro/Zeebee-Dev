import React, { useState } from 'react';
import { Button, List, Space } from 'antd';
import { DeleteOutlined, FileTextOutlined, UpOutlined } from '@ant-design/icons';
import ViewSvg from '@/assets/images/view.svg';
import styles from '../index.module.scss';
const DocumentItem = ({ columns, item, isHide }) => {
	const [expanded, setExpanded] = useState(false);
	return (
		<div className="border border-base-border bg-[#f5f5f5] rounded-md py-2 px-4 mb-4">
			<div className="flex gap-4">
				{columns.map((columnsItem, index) => (
					<div
						key={index}
						style={{
							width: columnsItem.width,
							minWidth: columnsItem.minWidth,
							flex: columnsItem.minWidth ? 1 : '0 0 auto',
						}}
						className={[
							columnsItem?.ellipsis && styles.ellipsis,
							index === columns.length - 1 && 'flex items-center',
							columnsItem?.align === 'right' && 'justify-end',
							columnsItem?.align === 'center' && 'justify-center',
						].join(' ')}
					>
						{item?.children?.length > 0 && !index && (
							<UpOutlined
								onClick={() => setExpanded(!expanded)}
								style={{
									transform: expanded ? 'rotate(180deg)' : 'rotate(0deg)',
									transition: 'transform 0.3s',
									fontSize: '12px',
									marginRight: '4px',
								}}
							/>
						)}
						{columnsItem?.render?.(item[columnsItem.dataIndex], item) ||
							item[columnsItem.dataIndex] ||
							'--'}
					</div>
				))}
			</div>
			{expanded && (
				<List
					itemLayout="horizontal"
					dataSource={item?.children}
					renderItem={(item, index) => (
						<List.Item key={index}>
							<div className="h-8 px-4 bg-white border border-base-border rounded-md flex items-center justify-center">
								<FileTextOutlined />
								<div className="ml-1">{item.title}</div>
							</div>
							<Space>
								<Button type="link" icon={<ViewSvg />}>
									view
								</Button>
								{!isHide && (
									<Button type="link" icon={<DeleteOutlined />} size="default">
										delete
									</Button>
								)}
							</Space>
						</List.Item>
					)}
				/>
			)}
		</div>
	);
};

export default DocumentItem;
